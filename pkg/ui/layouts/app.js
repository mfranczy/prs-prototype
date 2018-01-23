Vue.component('label-modal', {
    template: "#label-modal",
    data: function() {
        return {
            name: "",
            summary: "",
            fields: []
        }
    },
    methods: {
        clear: function() {
            this.name = "";
            this.summary = "";
            this.fields = [];
        },
        addField: function() {
            this.fields.push({name: "", description: ""});
        },
        removeField: function(index) {
            this.fields.splice(index, 1);
        },
        createLabel: function() {
            axios.post('/api/v1/labels/add', {
                name: this.name,
                summary: this.summary,
                fields: this.fields
            });
        }
    }
})

const Packages = {
    template: '#packages',
    data: function() {
        return {
            packages: this.fetchPackages(),
            labels: [],
            labelFields: [],
            attachPkgId: null,
            attachLabelId: null,
            fields: [ 'id', 'name', 'version', 'maintainer', 'labels', 'action' ],
            showModal: false
        }
    },
    methods: {
        fetchPackages: function() {
            axios.get('/api/v1/packages/list').then(response => { this.packages = response.data })
        },
        fetchLabels: function() {
            axios.get('/api/v1/labels/list').then(response => { this.labels = response.data })
        },
        attachLabel: function(id) {
            this.labelFields = [];
            this.attachPkgId = id;
            this.showModal = !this.showModal;
            this.labels = this.fetchLabels();
        },
        getLabel: function(id) {
            this.attachLabelId = id;
            this.labelFields = [];
            axios.get('/api/v1/labels/get/' + id).then(response => {
                let data = response.data[0];
                let ids = JSON.parse(data.fields_ids);
                let fields = JSON.parse(data.fields_names);
                let desc = JSON.parse(data.fields_descriptions);
                for (i in ids) {
                    if (ids[i] !== null) {
                        this.labelFields.push({
                            "id": ids[i],
                            "name": fields[i],
                            "desc": desc[i],
                            "value": null
                        })
                    }
                }
            })
        },
        attachReq: function() {
            let fields = [];
            for (i in this.labelFields) {
                if (this.labelFields[i].value !== null) {
                    fields.push({
                        "id": this.labelFields[i].id,
                        "value": this.labelFields[i].value
                    })
                }
            }
            axios.post('/api/v1/labels/attach-pkg', {
                "pkg_id": this.attachPkgId,
                "label_id": this.attachLabelId,
                "fields": fields
            });
        }
    },
    delimiters: ['${', '}']
}

const Labels = {
    template: '#labels',
    data: function() {
        return {
            labels: this.fetchLabels(),
            fields: [],
            labelFields: [],
            packages: null
        }
    },
    methods: {
        fetchLabels: function() {
            axios.get('/api/v1/labels/list').then(response => {
                this.labels = response.data;
                if (response.data)
                    this.fetchLabelPackages(response.data[0].id);
            })
        },
        fetchLabelPackages: function(labelId) {
            axios.all([
                axios.get('/api/v1/labels/get/' + labelId),
                axios.get('/api/v1/labels/list-packages/' + labelId)
            ]).then(response => {
                this.fields = [ 'name', 'version' ]
                if (response[0].data) {
                    this.labelFields = JSON.parse(response[0].data[0].fields_names)
                    this.fields = this.fields.concat(this.labelFields);
                }
                this.packages = response[1].data;
            });
        }
    },
    delimiters: ['${', '}']
}

const routes = [
  { path: '/packages', name: 'packages', component: Packages },
  { path: '/labels', name: 'labels', component: Labels },
  { path: '*', redirect: { name: 'packages' } }
]

const router = new VueRouter({
  routes
})

const app = new Vue({
    data: {
        showLabelModal: false
    },
    router
}).$mount('#app')
