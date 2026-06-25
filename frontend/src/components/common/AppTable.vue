<script setup>
defineProps({
  columns: {
    type: Array,
    required: true
  },
  data: {
    type: Array,
    required: true
  },
  onRowClick: {
    type: Function,
    default: null
  }
})
</script>

<template>
  <div class="table-container">
    <table class="table">
      <thead>
        <tr>
          <th v-for="(col, index) in columns" :key="index" :style="{ width: col.width }">
            {{ col.header }}
          </th>
        </tr>
      </thead>
      <tbody>
        <tr 
          v-for="(row, rowIndex) in data" 
          :key="row.id || rowIndex"
          @click="onRowClick ? onRowClick(row) : null"
          :style="{ cursor: onRowClick ? 'pointer' : 'default' }"
        >
          <td v-for="(col, colIndex) in columns" :key="colIndex">
            <!-- If column has a custom slot/render function, the parent will use scoped slots. 
                 Vue handles this differently from React. We use dynamic slots here. -->
            <slot :name="col.key" :row="row" :value="row[col.key]">
              {{ row[col.key] }}
            </slot>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
