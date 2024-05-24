<template>
  <v-data-table
    :headers="headers"
    :items="products"
    :items-per-page="50"
    class="products"
    @click:row="showDetail"
  >
    <template v-slot:top>
      <v-toolbar flat>
        <v-toolbar-title>{{ $t("product.products") }}</v-toolbar-title>
        <v-divider class="mx-4" inset vertical></v-divider>
        <v-spacer></v-spacer>
        <v-dialog
          v-if="role !== 'Viewer'"
          v-model="importDialog"
          max-width="500px"
        >
          <template v-slot:activator="{ on, attrs }">
            <Help
              title="Products management"
              keyString="help-productsManagement"
            />
            <v-btn
              color="primary"
              dark
              class="mb-2 mr-2"
              v-bind="attrs"
              v-on="on"
            >
              {{ $t("product.import.add") }}
            </v-btn>
          </template>
          <v-card>
            <v-form ref="form" v-model="validImport" lazy-validation>
              <v-card-title>
                <span class="text-h5">{{ $t("product.import.title") }}</span>
              </v-card-title>
              <v-card-text>
                <v-container>
                  <v-row>
                    <v-col cols="12">
                      <v-file-input
                        show-size
                        chips
                        accept=".xlsx"
                        v-model="file"
                        :label="$t('product.import.label')"
                      ></v-file-input>
                    </v-col>
                  </v-row>
                </v-container>
              </v-card-text>
              <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn color="blue darken-1" text @click="closeImport">
                  Cancel
                </v-btn>
                <v-btn
                  color="primary"
                  dark
                  :disabled="!validImport"
                  @click="importData"
                >
                  Import
                </v-btn>
              </v-card-actions>
            </v-form>
          </v-card>
        </v-dialog>

        <v-dialog v-model="dialog" max-width="1000px">
          <v-card>
            <v-card-title>
              <span class="text-h5"
                >#{{ detail.ProductCode }}, {{ $t("product.number") }}:
                {{ detail.Count }}, {{ $t("product.amount") }}: {{ detail.Avg }}
              </span>
            </v-card-title>

            <v-card-text>
              <v-container>
                <v-row>
                  <v-col cols="12">
                    <v-card class="v-card-home">
                      <v-card-title>{{ $t("product.orders") }}</v-card-title>
                      <apexchart
                        ref="realtimeOrd"
                        height="300"
                        type="area"
                        :options="barOptions"
                        :series="seriesOrd"
                      ></apexchart>
                    </v-card>
                  </v-col>
                  <v-col cols="12">
                    <v-data-table
                      :headers="headersDetail"
                      :items="productsOrders"
                      :items-per-page="10"
                      class="elevation-1"
                    >
                      <template v-slot:top>
                        <v-toolbar flat>
                          <v-toolbar-title>{{
                            $t("product.table.orders")
                          }}</v-toolbar-title>
                          <v-divider class="mx-4" inset vertical></v-divider>
                          <v-spacer></v-spacer>
                        </v-toolbar>
                      </template>
                      <template v-slot:no-data> </template>
                    </v-data-table>
                  </v-col>
                </v-row>
              </v-container>
            </v-card-text>

            <v-card-actions>
              <v-spacer></v-spacer>
              <v-btn color="blue darken-1" text @click="close"> Close </v-btn>
            </v-card-actions>
          </v-card>
        </v-dialog>
      </v-toolbar>
    </template>
    <template v-slot:no-data> </template>
  </v-data-table>
</template>
<script lang="ts">
import { Component, Vue, Watch } from "nuxt-property-decorator";
import Product from "~/model/Product";
import InfluxService from "~/services/InfluxService";
import { namespace } from "vuex-class";
const snackbar = namespace("Snackbar");
const spinner = namespace("Spinner");
import Help from "~/components/Help.vue";

@Component({
  components: { Help },
})
export default class ProductsPage extends Vue {
  @snackbar.Action
  public updateText!: (newText: string) => void;

  @snackbar.Action
  public updateColor!: (newColor: string) => void;

  @snackbar.Action
  public updateShow!: (newShow: boolean) => void;

  @spinner.Action
  public toggleSpinner!: (newShow: boolean) => void;

  $i18n: any;
  $t: any;
  products: Product[] = [];
  productsOrders: Array<any> = [];
  dialog: boolean = false;
  importDialog: boolean = false;
  detail: Product = {
    Id: "",
    ProductCode: "",
    Name: "",
    Avg: 0,
    Count: 0,
    Quantity: 0,
  };
  // @ts-ignore
  file: File = {};
  headers: Array<any> = [
    {
      // @ts-ignore
      text: this.$i18n.t("product.code"),
      align: "start",
      sortable: false,
      value: "ProductCode",
    },
    // @ts-ignore
    { text: this.$i18n.t("product.name"), value: "Name" },
    // @ts-ignore
    { text: this.$i18n.t("product.number"), value: "Count" },
    // @ts-ignore
    { text: this.$i18n.t("product.amount"), value: "Avg" },
    // @ts-ignore
    { text: this.$i18n.t("product.quantity"), value: "Quantity" },
  ];
  headersDetail: Array<any> = [
    {
      // @ts-ignore
      text: this.$i18n.t("product.created"),
      align: "start",
      sortable: false,
      value: "CreatedAt",
    },
    // @ts-ignore
    { text: this.$i18n.t("product.number"), value: "Count" },
    // @ts-ignore
    { text: this.$i18n.t("product.amount"), value: "Avg" },
  ];
  title: string = "Products";
  barOptions: any = {
    chart: {
      id: "product-detail",
    },
    tools: {
      download: true,
      selection: true,
      zoom: true,
      zoomin: true,
      zoomout: true,
      pan: true,
      reset: true,
    },
    xaxis: {
      type: "datetime",
      tickAmount: 6,
    },
    dataLabels: {
      enabled: false,
    },
    markers: {
      size: 0,
      style: "hollow",
    },
    stroke: {
      curve: "smooth",
    },
    annotations: {
      xaxis: [
        {
          x: new Date().getTime(),
          borderColor: "#999",
          yAxisIndex: 0,
          label: {
            show: true,
            text: "Today",
            style: {
              color: "#fff",
              background: "#775DD0",
            },
          },
        },
      ],
    },
  };
  seriesOrd: Array<any> = [
    {
      // @ts-ignore
      name: this.$i18n.t("dash.orders"),
      data: [],
    },
    {
      // @ts-ignore
      name: this.$i18n.t("dash.realOrders"),
      data: [],
    },
  ];
  $axios: any;
  ord: Array<any> = [];
  histOrd: Array<any> = [];
  $refs: any;
  validImport: boolean = false;
  role: string = "";
  $auth: any;

  head() {
    return {
      title: this.title,
    };
  }

  @Watch("dialog")
  onPropertyChanged(value: string, oldValue: string) {
    value || this.close();
  }

  @Watch("$route.query.storeId")
  onPropertyChangedDelete(value: string, oldValue: string) {
    this.getProducts(value);
  }

  mounted() {
    if (this.$route.query.storeId !== undefined) {
      this.getProducts(this.$route.query.storeId);
    }
    this.role = this.$auth.$state.user.Role;
  }

  getProducts(storeId: any) {
    this.$axios
      .get("/products/" + storeId + "/1000/0")
      .then((response: any) => {
        if (response.data.success === true) {
          if (response.data.products !== null) {
            response.data.products.forEach((item: Product) => {
              this.products.push({
                Id: item.Id,
                ProductCode: item.ProductCode,
                Name: item.Name,
                Avg: Math.ceil(item.Avg),
                Count: item.Count,
                Quantity: item.Quantity,
              });
            });
          }
        } else {
          this.updateText(response.data.error);
          this.updateColor("red");
          this.updateShow(true);
        }
      });
  }

  showDetail(item: Product, storeId: string) {
    this.dialog = true;
    this.getProductDetailData(item);
    this.getOrderHistory(item, storeId);
  }

  close() {
    this.dialog = false;
    this.histOrd = [];
    this.ord = [];
    this.$refs.realtimeOrd.updateSeries([
      { data: this.ord },
      { data: this.histOrd },
    ]);
  }

  closeImport() {
    this.importDialog = false;
    this.$refs.form.resetValidation();
  }

  getProductDetailData(item: Product) {
    this.detail = item;
    this.histOrd = [];
    this.ord = [];
    this.$axios
      .get(
        "/prediction/query/products/" +
          item.ProductCode +
          "/" +
          this.$route.query.storeId +
          "/1/1"
      )
      .then((response: any) => {
        if (response.data.success === true) {
          if (response.data.products !== null) {
            this.histOrd = response.data.prediction
              .filter((item: any) => item.Index === "d0")
              .map((record: any) => [record.Date.split(" ")[0], record.Val]);
            for (let d = 0; d < 60; d++) {
              let date = new Date();
              date.setDate(date.getDate() + d);
              const day = date.toISOString().slice(0, 10);
              const base = response.data.prediction.filter(
                (item: any) => item.Date.split(" ")[0] === day
              );
              if (base.length > 0) {
                this.ord.push([date.getTime(), base[0].Val]);
              }
            }
            setTimeout(() => {
              this.$refs.realtimeOrd.updateSeries([
                { data: this.ord },
                { data: this.histOrd },
              ]);
            }, 500);
          }
        } else {
          this.updateText(response.data.error);
          this.updateColor("red");
          this.updateShow(true);
        }
      });
  }

  getOrderHistory(item: Product, storeId: string) {
    this.productsOrders = [];
    this.$axios
      .get(
        this.$config.internalApi +
          "orders/product/" +
          storeId +
          "/" +
          item.ProductCode
      )
      .then((response: any) => {
        if (response.data.success === true) {
          this.productsOrders = response.data.orders;
        }
      });
  }

  importData() {
    this.toggleSpinner(true);
    let formData = new FormData();
    formData.append("importFile", this.file, this.file.name);
    this.$axios
      .post("/warehouse/import/" + this.$route.query.storeId, formData)
      .then((response: any) => {
        if (response.data.success === true) {
          this.$refs.form.resetValidation();
          this.updateText(this.$i18n.t("product.import.success"));
          this.updateColor("green");
          this.updateShow(true);
          this.closeImport();
          this.toggleSpinner(false);
        } else {
          this.$refs.form.resetValidation();
          this.updateText(response.data.error);
          this.updateColor("red");
          this.updateShow(true);
          this.toggleSpinner(false);
        }
      });
  }
}
</script>

<style scoped>
td {
  cursor: pointer !important;
}
.v-dialog .v-card__title {
  background: #4633af !important;
  color: #fff !important;
}
.v-dialog .container {
  color: #000 !important;
}
.theme--dark .v-dialog .v-card__title {
  background-color: #4f5249 !important;
}
</style>
