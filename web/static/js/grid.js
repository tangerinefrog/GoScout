let gridApi = null;

async function renderGrid() {
    const rowData = await getRows();
    const theme = agGrid.themeAlpine.withPart(agGrid.colorSchemeDark);

    const gridOptions = {
        rowData: rowData,
        pagination: true,
        paginationPageSize: 50,
        paginationPageSizeSelector: [50, 100, 200],
        enableCellTextSelection: true,
        columnDefs: defineColumns(),
        theme: theme,
    };

    const gridElem = $('#jobsGrid')[0];
    gridApi = agGrid.createGrid(gridElem, gridOptions);
}

async function getRows() {
    const jobs = await getJobs();

    const rows = jobs.map(j => {
        const date = new Date(j.date_posted);
        const dateFormatted = dayjs(date).format('MMM DD HH:mm');

        return {
            id: j.id,
            title: { title: j.title, url: j.url },
            company: j.company,
            location: j.location,
            status: j.status,
            grade: { grade: j.grade, grade_reasoning: j.grade_reasoning},
            date_posted: dateFormatted,
            note: j.note,
        }
    });

    return rows;
}

async function refreshGrid() {
    const rows = await getRows();
    gridApi.setGridOption("rowData", rows);
}

function defineColumns() {
    return [
        {
            field: 'date_posted',
            headerName: 'Date',
            type: 'dateTime',
            width: 130,
        },
        {
            field: 'title',
            width: 300,
            cellRenderer: titleRenderer,
        },
        {
            field: 'company',
            width: 200,
        },
        {
            field: 'location',
            width: 200,
            filter: true,
            filterParams: getFilterParams(),
        },
        {
            field: 'grade',
            width: 80,
            cellRenderer: gradeRenderer,
            comparator: gradeSorter
        },
        {
            field: 'status',
            width: 100,
            editable: true,
            cellEditor: 'agSelectCellEditor',
            cellEditorParams: {
                values: statuses,
            },
            onCellValueChanged: onEdit,
            filter: true,
            filterParams: getFilterParams(),
            cellStyle: statusCellStyle
        },
        {
            field: 'note',
            flex: 1,
            editable: true,
            cellEditor: 'agTextCellEditor',
            onCellValueChanged: onEdit,
        },
    ];
}

function titleRenderer(cell) {
    return `<a href="${cell.value.url}" target="_blank">${cell.value.title}</a>`;
}

function gradeRenderer(cell) {
    const value = cell.value.grade;
    const cellDiv = document.createElement('div');
    cellDiv.innerText = value;
    cellDiv.style.cursor = 'pointer';
    cellDiv.style.textAlign = 'center';

    cellDiv.addEventListener('click', () => {
        $('.js-grading-modal-text').text(cell.value.grade_reasoning);
        MicroModal.show('grading-modal');
    });

    return cellDiv;
}

function gradeSorter(valueA, valueB) {
    if (valueA.grade == valueB.grade) {
        return 0;
    }
    
    return (valueA.grade > valueB.grade) ? 1 : -1;
} 

function statusCellStyle(cell) {
    const style = {
        textAlign: 'center',
    }

    switch (cell.value) {
        case 'created':
            style.color = '#ffffff';
            style.backgroundColor = '#2b2b2b';
            break;
        case 'graded':
            style.color = '#000000';
            style.backgroundColor = '#efc004f0';
            style.borderRadius = '8px';
            break;
        case 'ignored':
            style.color = '#ffffffff';
            style.backgroundColor = '#a60e0ea9';
            style.borderRadius = '8px';
            break;
        case 'applied':
            style.color = '#ffffffff';
            style.backgroundColor = '#1c8305ff';
            style.borderRadius = '8px';
            break;
    }

    return style;
}

function getFilterParams() {
    return {
        closeOnApply: true,
        filterOptions: ['contains', 'notContains'],
    };
}

async function onEdit(e) {
    const id = e.data.id;
    const fieldName = e.colDef.field;
    let value = e.newValue;

    if (e.colDef.cellDataType === 'text' && value === null) {
        value = '';
    }

    const isUpdated = await updateJob(id, fieldName, value);
    if (!isUpdated) {
        e.data[fieldName] = e.oldValue;
        e.api.refreshCells();
        showErrorToast('Edit failed due to server error');
    }
}