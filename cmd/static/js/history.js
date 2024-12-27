if (document.getElementById("default-table") && typeof simpleDatatables.DataTable !== 'undefined') {
    const dataTable = new simpleDatatables.DataTable("#default-table", {
        paging: true,
        perPage: 250,
        perPageSelect: [13, 50, 100, 250, 500, 1000, 2000, 4000, 8000],
        searchable: true,
        sortable: false,
        classes: {
            dropdown: "datatable-dropdown mr-5",
            input: "datatable-input ml-5",
            top: "datatable-top mt-5",
            // add custom HTML classes, full list: https://fiduswriter.github.io/simple-datatables/documentation/classes
            // we recommend keeping the default ones in addition to whatever you want to add because Flowbite hooks to the default classes for styles
        },
    })
}