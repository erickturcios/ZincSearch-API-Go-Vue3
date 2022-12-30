<script setup>
import { ref } from 'vue'


const endpoint = "http://localhost:8001";
const records = ref(null);
const currentText = ref("");
const currentSelected = ref(0);

const query = ref("");

//verifica si la linea actual ha sido seleccionada
function isSelected(index){
  return currentSelected.value != null && index == currentSelected.value;
}

function select(index){
  currentSelected.value = index;
  if(typeof records != "undefined"
    && typeof records.value[index] != "undefined"
    && typeof records.value[index]._source != "undefined"
    && typeof records.value[index]._source.TextBody != "undefined"
  ){
    currentText.value = records.value[index]._source.TextBody;
  }
}

async function receiveRecords(response){
  let data = await response.json();

  if(typeof data.hits != 'undefined' 
     && typeof data.hits.total != 'undefined'
     && typeof data.hits.total.value != 'undefined'
     && typeof data.hits.hits != 'undefined'
     && data.hits.total.value > 0
  ){
    records.value = data.hits.hits;
    select(0);
  }
  else{
    records.value = [];
  }
}


function searchRecords(){  
  records.value = [];
  currentText.value = "";
  fetch(endpoint + `/search?query=` + query.value)
    .then((response) => receiveRecords(response))
}
</script>

<template>
  <div class="container" style="padding:30px;">
    <div class="row justify-content-md-center">
      
      <div class="col-sm-6">
        <div class="input-group mb-3">
          <input type="text" class="form-control" placeholder="Escriba texto de busqueda" aria-label="Texto"
            aria-describedby="basic-addon2" v-model="query" @keyup.enter="searchRecords">
          <div class="input-group-append">
            <button class="btn btn-success" type="button" @click="searchRecords">Buscar</button>
          </div>
        </div>
      </div>
    </div>
    <div class="row" >
      <div class="col-sm-6 table-responsive">
          <table class="table table-hover table-bordered border-secondary table-sm"  >
              <thead>
                  <tr>
                    <th scope="col">#</th>
                    <th scope="col">De</th>
                    <th scope="col">Para</th>
                    <th scope="col">Titulo</th>
                  </tr>
              </thead>
              <tbody>
                <tr v-for="(item, index) in records" 
                    v-bind:class="{'table-active':isSelected(index)}"
                    @click="select(index)"
                    >
                  <td>{{index+1}}</td>
                  <td>{{item._source.From}}</td>
                  <td>{{item._source.To}}</td>
                  <td>{{item._source.Subject}}</td>
                </tr>

              </tbody>
          </table>
      </div>
      <div class="col-sm-5">
        <div class="container border rounded border-4 justify-content-md-center text-wrap text-break">
          <pre><span v-html="currentText" style="white-space: pre;"></span></pre>
        </div>
        
      </div>
    </div>
  </div>
</template>
