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
            fields: [ 'id', 'name', 'version', 'maintainer', 'labels', 'action' ]
        }
    },
    methods: {
        fetchPackages: function() {
            axios.get('/api/v1/packages/list').then(response => { this.packages = response.data })
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
