if (document.getElementById("default-table") && typeof simpleDatatables.DataTable !== 'undefined') {
    const dataTable = new simpleDatatables.DataTable("#default-table", {
        paging: true,
        perPage: 250,
        perPageSelect: [13, 50, 100, 250, 500, 1000, 2000, 4000, 8000],
        searchable: true,
        sortable: false,
        scrollY: "74vh",
        classes: {
            active: "datatable-active",
            bottom: "datatable-bottom",
            container: "datatable-container",
            cursor: "datatable-cursor",
            dropdown: "datatable-dropdown mr-5",
            ellipsis: "datatable-ellipsis",
            empty: "datatable-empty",
            headercontainer: "datatable-headercontainer",
            info: "datatable-info",
            input: "datatable-input ml-5",
            loading: "datatable-loading",
            pagination: "datatable-pagination",
            paginationList: "datatable-pagination-list",
            search: "datatable-search",
            selector: "datatable-selector",
            sorter: "datatable-sorter",
            table: "datatable-table w-full text-sm text-left rtl:text-right text-gray-100 dark:text-gray-100",
            top: "datatable-top mt-5",
            wrapper: "datatable-wrapper"
            // add custom HTML classes, full list: https://fiduswriter.github.io/simple-datatables/documentation/classes
            // we recommend keeping the default ones in addition to whatever you want to add because Flowbite hooks to the default classes for styles
        },
    })
}