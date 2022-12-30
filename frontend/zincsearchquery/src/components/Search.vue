<script setup>
import { ref } from 'vue'


const endpoint = "http://localhost:8001";
const records = ref(null);
const currentText = ref("");
const currentSelected = ref(0);
//paginacion
const pagesPrevious = ref(0)
const pagesTotal = ref(1);
const pageCurrent = ref(1);
const disableBack = ref(true)
const disableForward = ref(true)

const query = ref("");

//verifica si la linea actual ha sido seleccionada
function isSelected(index) {
  return currentSelected.value != null && index == currentSelected.value;
}

function select(index) {
  currentSelected.value = index;
  if (typeof records != "undefined"
    && typeof records.value[index] != "undefined"
    && typeof records.value[index]._source != "undefined"
    && typeof records.value[index]._source.TextBody != "undefined"
  ) {
    currentText.value = records.value[index]._source.TextBody;
  }
}

async function receiveRecords(response, pageIdx) {
  let data = await response.json();

  pageCurrent.value = pageIdx;
  pagesPrevious.value = 0
  disableBack.value = true
  disableForward.value = true

  if (pageCurrent.value > 1) {
    pagesPrevious.value = (pageCurrent.value - 1) * 10
  }

  if (typeof data.hits != 'undefined'
    && typeof data.hits.total != 'undefined'
    && typeof data.hits.total.value != 'undefined'
    && typeof data.hits.hits != 'undefined'
    && data.hits.total.value > 0
  ) {
    records.value = data.hits.hits;
    pagesTotal.value = Math.ceil(data.hits.total.value / 10)
    disableBack.value = pageCurrent.value <= 1
    disableForward.value = pageCurrent.value >= pagesTotal.value
    select(0);
  }
  else {
    records.value = [];
  }
}

//utilizado para paginacion de resultados
function searchRecordsForPage(pageIdx) {
  if (pageIdx < 1 || pageIdx > pagesTotal.value) {
    return;
  }
  records.value = [];
  currentText.value = "";
  fetch(endpoint + `/search?query=` + query.value + `&page=` + pageIdx)
    .then((response) => receiveRecords(response, pageIdx))
}

searchRecordsForPage(1);

</script>

<template>
  <div class="container" style="padding:30px;">
    <div class="row justify-content-md-center">

      <div class="col-sm-6">
        <div class="input-group mb-3">
          <input type="text" class="form-control" placeholder="Escriba texto de busqueda" aria-label="Texto"
            aria-describedby="basic-addon2" v-model="query" @keyup.enter="searchRecordsForPage(1)">
          <div class="input-group-append">
            <button class="btn btn-success" type="button" @click="searchRecordsForPage(1)">Buscar</button>
          </div>
        </div>
      </div>
    </div>
    <div class="row">
      <!-- Navegacion -->
      <nav aria-label="Navegacion">
        <ul class="pagination justify-content-center">
          <li class="page-item" v-bind:class="{ 'disabled': disableBack }">
            <a class="page-link" href="#" tabindex="-1"
              @click.prevent="searchRecordsForPage(pageCurrent - 1)">Anterior</a>
          </li>
          <li class="page-item disabled">
            <a class="page-link" href="#">{{ pageCurrent }} <span class="sr-only">( de {{ pagesTotal }})</span></a>
          </li>
          <li class="page-item" v-bind:class="{ 'disabled': disableForward }">
            <a class="page-link" href="#" @click.prevent="searchRecordsForPage(pageCurrent + 1)">Siguiente</a>
          </li>
        </ul>
      </nav>
      <!-- Listado -->
      <div class="col-sm-6 table-responsive">
        <table class="table table-hover table-bordered border-secondary table-sm">
          <thead>
            <tr>
              <th scope="col">#</th>
              <th scope="col">De</th>
              <th scope="col">Para</th>
              <th scope="col">Titulo</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(item, index) in records" v-bind:class="{ 'table-active': isSelected(index) }"
              @click="select(index)">
              <td>{{ pagesPrevious + index + 1 }}</td>
              <td>{{ item._source.From }}</td>
              <td>{{ item._source.To }}</td>
              <td>{{ item._source.Subject }}</td>
            </tr>

          </tbody>
        </table>
      </div>

      <!-- Contenido registro seleccionado -->
      <div class="col-sm-5">
        <div class="container border rounded border-4 justify-content-md-center text-wrap text-break">
          <pre><span v-html="currentText" style="white-space: pre;"></span></pre>
        </div>

      </div>
    </div>
  </div>
</template>
