async function setupGrid() {
    const rowData = await getRows();
    const theme = agGrid.themeAlpine.withPart(agGrid.colorSchemeDark)

    const gridOptions = {   
        rowData: rowData,
        pagination: true,
        paginationPageSize: 50,
        paginationPageSizeSelector: [50, 100, 200],
        enableCellTextSelection: true,
        autoSizeStrategy: {
            defaultMinWidth: 100,
        },
        
        columnDefs: [
            { field: 'date_posted', headerName: 'Date', type: 'dateTime', width: 120 },
            { field: 'title', width: 300, cellRenderer: titleRenderer },
            { field: 'company', width: 200 },
            { field: 'location', width: 200, filter: true, filterParams: getFilterParams() },
            { field: 'num_applicants', headerName: 'Applicants',  width: 150, sortable: false },
            { field: 'status', width: 120, editable: true, cellEditor: 'agSelectCellEditor', cellEditorParams: { values: statuses, }, onCellValueChanged: onEdit, filter: true, filterParams: getFilterParams() },
            { field: 'note', flex: 1,  editable: true, cellEditor: 'agTextCellEditor', onCellValueChanged: onEdit },
        ],
        
        theme: theme,
    };

    const myGridElement = document.querySelector('#jobsGrid');
    agGrid.createGrid(myGridElement, gridOptions);
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
            num_applicants: j.num_applicants,
            date_posted: dateFormatted,
            note: j.note,
        }
    });

    return rows;
}

function titleRenderer(cell) {
    return `<a href="${cell.value.url}" target="_blank">${cell.value.title}</a>`;
}

function getFilterParams() {
    return {
        closeOnApply: true,
        filterOptions: ['contains', 'notContains']
    }
}

async function onEdit(e) {
    const id = e.data.id;
    const fieldName = e.colDef.field;
    let value = e.newValue;

    if(e.colDef.cellDataType === 'text' && value === null) {
        value = '';
    }

    const isUpdated = await updateJob(id, fieldName, value);
    if (!isUpdated) {
        e.data[fieldName] = e.oldValue;
        e.api.refreshCells();
        showErrorToast('Edit failed due to server error');
    }
}

setupGrid();
