<template>
    <div>
<!--        <div>route.name: {{ $route.name }}</div>-->
<!--        <div>route.path: {{ $route.path }}</div>-->
<!--        <div>route.params: {{ $route.params }}</div>-->
<!--        <div>route.query: {{ $route.query }}</div>-->
<!--        <div>route.hash: {{ $route.hash }}</div>-->
<!--        <div>route.fullPath: {{ $route.fullPath }}</div>-->
        
        <table border="1">
        <tr>
            <th>Name</th>
            <th>Size</th>
            <th>ModTime</th>
            <th>Actions</th>
        </tr>
        <tr v-for="item in items">
            <td>{{ item.name }}</td>
            <td>{{ item.size }}</td>
            <td>{{ item.modTime }}</td>
            <td><a v-if="item.type == 'file'" :href="'/-/files/' + item.path + '?action=download'">Download</a></td>
        </tr>
    </table>
    </div>
</template>

<script>
import axios from 'axios'

export default {
    name: 'files',
    components: {
    },
    mounted (){
        let pathMatch = this.$route.params.pathMatch
        axios.get('/-/files/' + pathMatch)
        .then( response => {
            console.log(response)
            this.items = response.data
        }).catch(function (error) {
            console.log(error);
     });
        console.log('Files mounted')
    },
    watch: {
        '$route' (to, from){
            console.log("route changed from {} to {}", from, to)
            let pathMatch = to.params.pathMatch;
            axios.get("/-/files/" + pathMatch)
                .then( response => {
                    console.log(response)
                    this.items = response.data
                }).catch(function (error) {
                console.log(error);
            });
        }
    },
    data () {
        return {
            items: [
                {
                    name: 'hello.txt',
                    size: '10kb',
                    modTime: '2019-10-01 10:10:00',
                    actions: 'Download,View'
                },
                {
                    name: 'world.txt',
                    size: '15kb',
                    modTime: '2019-08-10 12:25:00',
                    actions: 'Download,View'
                }
            ]
        }
    }
}
</script>