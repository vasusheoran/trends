if (document.getElementById("default-table") && typeof simpleDatatables.DataTable !== 'undefined') {
    const dataTable = new simpleDatatables.DataTable("#default-table", {
        paging: true,
        perPage: 250,
        perPageSelect: [50, 100, 250],
        searchable: true,
        sortable: false,
        show: false
    }).afterMain();
}