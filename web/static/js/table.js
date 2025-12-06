async function getRows() {
    let rows = [];

    const resp = await fetch(`/api/jobs`);
    if (resp.ok) {
        const jobs = await resp.json();
        rows = jobs.map(j => {
            const date = new Date(j.date_posted);
            const dateFormatted = dayjs(date).format('MMM DD HH:mm');
            
            return {
                id: j.id,
                title: j.title,
                company: j.company,
                status: j.status,
                num_applicants: j.num_applicants,
                date_posted: dateFormatted,
            }
        });
    }

    return rows;
}

async function setupGrid() {
    const rowData = await getRows();
    const theme = agGrid.themeAlpine.withPart(agGrid.colorSchemeDark)

    const gridOptions = {   
        rowData: rowData,
        pagination: true,
        paginationPageSize: 50,
        paginationPageSizeSelector: [50, 100, 200],
        columnDefs: [
            { field: 'title', width: 300 },
            { field: 'company', width: 200 },
            { field: 'num_applicants', width: 150 },
            { field: 'status', width: 120 },
            { field: 'date_posted', type: 'dateTime', width: 120 },
        ],

        theme: theme,
    };

    const myGridElement = document.querySelector('#jobsGrid');
    agGrid.createGrid(myGridElement, gridOptions);
}

setupGrid();
