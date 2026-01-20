renderGrid();

MicroModal.init({
    onShow: modal => {
        if (modal.id === 'options-modal') {
            fillConfigValues();
        }
    },
});

async function fillConfigValues() {
    const config = await getConfig();
    if (!config) {
        return;
    }

    $('#grading-requirements-textarea').val(config.grading_profile);
    $('#search-query-input').val(config.search_query);
    $('#search-filter-input').val(config.search_filter);
    $('#scraping-interval-input').val(config.search_period_hours);
}

async function onConfigSaveBtnClick() {
    const gradingRequirements = $('#grading-requirements-textarea').val();
    const searchQuery = $('#search-query-input').val();
    const searchFilter = $('#search-filter-input').val();
    const scrapingInterval = Number($('#scraping-interval-input').val());

    const isSaved = await updateConfig({
        gradingRequirements: gradingRequirements,
        searchQuery: searchQuery,
        searchFilter: searchFilter,
        scrapingInterval: scrapingInterval
    });

    if (isSaved) {
        showSuccessToast('Config saved successfully');
        MicroModal.close('options-modal');
    } else {
        showErrorToast('Config save failed');
    }
}

async function onScrapeBtnClick() {
    const ok = await scrapeJobs();
    if (!ok) {
        showErrorToast('Jobs scraping failed');
    } else {
        showSuccessToast('Jobs scraped successfully');
        await refreshGrid();
    }
}

async function onStartGradingBtnClick() {
    const ok = await startGrading();
    if (!ok) {
        showErrorToast('Grading start failed');
    } else {
        showSuccessToast('Grading started');
    }
}

async function onGradeStatusBtnClick() {
    const status = await getGradingStatus();
    if (!status) {
        showErrorToast('Failed to get grading status');
    } else {
        $('.js-grade-status-modal-text').text(status);
        MicroModal.show('grade-status-modal');
    }
}

async function onSearchBtnClick() {
    await refreshGrid();
}