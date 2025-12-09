const statuses = [
    'created',
    'graded',
    'ignored',
    'applied',
]

async function getJobs() {
    const resp = await fetch(`/api/jobs`);
    if (!resp.ok) {
        showErrorToast('Jobs fetch failed');
        return [];
    }

    const jobs = await resp.json();
    if (!jobs) {
        return [];
    }

    return jobs;
}

async function updateJob(id, field, value) {
    const reqBody = {
        [field]: value
    };

    const resp = await fetch(`/api/jobs/${id}`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(reqBody),
    });

    return resp.ok;
}

async function getConfig() {
    const resp = await fetch(`/api/config`);
    if (!resp.ok) {
        showErrorToast('Config fetch failed');
        return null;
    }

    const config = await resp.json();
    return config;
}

async function updateConfig(config) {
    const reqBody = {
        grading_profile: config.gradingRequirements,
        search_query: config.searchQuery,
        search_filter: config.searchFilter,
        search_period_hours: config.scrapingInterval,
    };

    try {
        const resp = await fetch(`/api/config`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(reqBody),
        });
    
        if (!resp.ok) {
            const body = await resp.json();
            if (body) {
                console.error(`Error occured while saving the config: ${body}`);
            }
        }
    
        return resp.ok;
    }
    catch {
        return false;
    }
}

async function scrapeJobs() {
    preloader.show();
    try {
        const resp = await fetch(`/api/scrape`, {
            method: 'POST'
        });
    
        if (!resp.ok) {
            const body = await resp.json();
            if (body) {
                console.error(`Error occured while saving the config: ${body}`);
            }
        }
    
        preloader.hide();
        return resp.ok;
    }
    catch {
        preloader.hide();
        return false;
    }
}
