<template>
  <div>
    <v-row justify="center" align="center" v-if="waitingForData">
      <v-col cols="12">
        <v-card class="v-card-home">
          <v-card-title>
            {{ $t("index.income") }}
            <v-spacer />
            <Help title="Income" keyString="help-income" />
          </v-card-title>
          <apexchart
            ref="realtimeInc"
            height="300"
            type="area"
            :options="lineOptions"
            :series="seriesInc"
          ></apexchart>
        </v-card>
      </v-col>
    </v-row>
    <v-row v-if="waitingForData">
      <v-col cols="12" sm="12" xs="12" md="6">
        <v-card class="v-card-home">
          <v-card-title>
            {{ $t("index.order") }}
            <v-spacer />
            <Help title="Orders" keyString="help-orders" />
          </v-card-title>
          <apexchart
            ref="realtimeOrd"
            height="300"
            type="bar"
            :options="barOptions"
            :series="seriesOrders"
          ></apexchart>
        </v-card>
      </v-col>
      <v-col cols="12" sm="12" xs="12" md="6">
        <v-card class="v-card-home">
          <v-card-title>
            {{ $t("index.visitors") }}
            <v-spacer />
            <Help title="Trends" keyString="help-trends" />
          </v-card-title>
          <apexchart
            ref="realtimeVisitors"
            height="300"
            type="line"
            :options="trendOptions"
            :series="seriesVisitors"
          ></apexchart>
        </v-card>
      </v-col>
    </v-row>
    <v-row v-if="waitingForData">
      <v-col cols="12">
        <v-data-table
          v-model="selected"
          :headers="toStoreHeaders"
          :items="products"
          :single-select="singleSelect"
          show-select
          item-key="ProductCode"
          class="elevation-1"
          :search="date"
          :custom-filter="filterByDate"
        >
          <template v-slot:top>
            <v-toolbar flat>
              <v-toolbar-title>{{ $t("index.products") }}</v-toolbar-title>
              <v-spacer></v-spacer>
              <v-menu
                v-model="menu"
                :close-on-content-click="false"
                :nudge-right="40"
                transition="scale-transition"
                offset-y
                min-width="auto"
              >
                <template v-slot:activator="{ on, attrs }">
                  <v-text-field
                    v-model="date"
                    :label="$t('index.dateToStore')"
                    prepend-icon="mdi-calendar"
                    readonly
                    v-bind="attrs"
                    v-on="on"
                    class="pt-5 mr-2 dateSearch"
                    hint="MM/DD/YYYY format"
                    persistent-hint
                  ></v-text-field>
                </template>
                <v-date-picker
                  v-model="date"
                  @input="menu = false"
                ></v-date-picker>
              </v-menu>
              <v-text-field
                    v-model="productFilter"
                    :label="$t('index.productFilter')"
                    class="pt-5 mr-2 productSearch"
                  ></v-text-field>
              <v-btn @click="clearDate()" class="mr-2" color="primary">
                {{ $t("index.clearDate") }}
              </v-btn>
              <v-btn
                v-if="role !== 'Viewer'"
                @click="makeOrder()"
                class="mr-2"
                color="primary"
                :disabled="selected.length === 0"
              >
                {{ $t("index.makeOrder") }}
              </v-btn>
              <Help
                title="Products to order"
                keyString="help-productsToOrder"
              />
            </v-toolbar>
          </template>
          <template v-slot:item.Id="{ item }">
            <v-tooltip bottom>
              <template v-slot:activator="{ on, attrs }">
                <v-btn
                  v-bind="attrs"
                  v-on="on"
                  @click="solved(item)"
                  class="mr-2"
                  color="primary"
                >
                  <v-icon small> mdi-check-bold </v-icon>
                </v-btn>
              </template>
              <span>{{ $t("tooltips.solved") }}</span>
            </v-tooltip>
          </template>
          <template v-slot:item.DateToNeed="{ item }">
            <span>{{ new Date(item.DateToNeed).toLocaleDateString() }}</span>
          </template>
          <template v-slot:item.DateToOrder="{ item }">
            <span>{{ new Date(item.DateToOrder).toLocaleDateString() }}</span>
          </template>
          <template v-slot:item.Diff="{ item }">
              <span v-if="item.Diff <= 0" class="redColor">{{  item.Diff }}</span>
              <span v-if="item.Diff > 0" class="greenColor">{{  item.Diff }}</span>
          </template>
        </v-data-table>
      </v-col>
    </v-row>
    <v-row justify="center" align="center" v-if="!waitingForData">
      <v-alert border="top" color="red lighten-2" dark>
        {{ $t("validation.noData") }}
      </v-alert>
    </v-row>
    <v-dialog v-model="makeOrderDialog" max-width="500px">
      <v-card>
        <v-card-title class="text-h5 white--text">{{
          $t("index.makeTitle")
        }}</v-card-title>
        <v-card-text>
          <v-select
            v-model="selectedSupplier"
            :items="suppliers"
            :label="$t('index.selectSupplier')"
            item-text="Name"
            item-value="Id"
            required
          ></v-select>
          <v-card-actions>
            <v-btn color="blue darken-1" text @click="closeMakeOrder">{{
              $t("index.cancel")
            }}</v-btn>
            <v-btn color="primary" @click="sendOrderToSupplier">{{
              $t("index.sendOrder")
            }}</v-btn>
          </v-card-actions>
          <v-spacer></v-spacer>
        </v-card-text>
      </v-card>
    </v-dialog>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue, Watch } from "nuxt-property-decorator";
import Product from "~/model/Product";
import Store from "~/model/Store";
import Supplier from "~/model/Supplier";
import Route from "~/model/Route";
import { namespace } from "vuex-class";
import Help from "~/components/Help.vue";
import ProductToStore from "~/model/ProductToStore";

const snackbar = namespace("Snackbar");
const spinner = namespace("Spinner");

@Component({
  components: { Help },
})
export default class IndexPage extends Vue {
  @Prop() readonly store!: Store;

  @snackbar.Action
  public updateText!: (newText: string) => void;

  @snackbar.Action
  public updateColor!: (newColor: string) => void;

  @snackbar.Action
  public updateShow!: (newShow: boolean) => void;

  @spinner.Action
  public toggleSpinner!: (newShow: boolean) => void;

  // important for nuxt ts
  $i18n: any;
  $refs: any;
  $nuxt: any;
  $auth: any;
  $t: any;
  $axios: any;

  singleSelect: boolean = false;
  selected: Array<any> = [];
  products: Array<any> = [];
  menu: boolean = false;
  calculating: boolean = false;
  date: string = "";
  waitingForData: boolean = true;
  productFilter: string = "";
  toStoreHeaders: Array<any> = [
    // @ts-ignore
    { text: this.$i18n.t("index.product"), value: "ProductCode" },
    // @ts-ignore
    { text: this.$i18n.t("index.name"), value: "Name" },
    // @ts-ignore
    { text: this.$i18n.t("index.expected"), value: "Expected" },
    // @ts-ignore
    { text: this.$i18n.t("index.real"), value: "Quantity" },
    // @ts-ignore
    { text: this.$i18n.t("index.diff"), value: "Diff" },
    // @ts-ignore
    { text: this.$i18n.t("index.dateToOrder"), value: "DateToOrder" },
    // @ts-ignore
    { text: this.$i18n.t("index.dateToStore"), value: "DateToNeed" },
    {
      // @ts-ignore
      text: "",
      align: "start",
      sortable: false,
      value: "Id",
    },
  ];
  value: Array<any> = [];
  ord: Array<any> = [];
  invoices: Array<any> = [];
  inc: Array<any> = [];
  histVisit: Array<any> = [];
  histOrd: Array<any> = [];
  histInc: Array<any> = [];
  lineOptions: any = {
    chart: {
      id: "order-line",
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
      width: 2,
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
  barOptions: any = {
    chart: {
      id: "order-bar",
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
      width: 2,
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
  trendOptions: any = {
    chart: {
      id: "trend-line",
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
      size: 5,
    },
    stroke: {
      curve: "smooth",
      width: 2,
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
  seriesOrders: Array<any> = [
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
  seriesVisitors: Array<any> = [
    {
      // @ts-ignore
      name: this.$i18n.t("dash.visit"),
      data: [],
    },
    {
      // @ts-ignore
      name: this.$i18n.t("dash.realVisit"),
      data: [],
    },
  ];
  seriesInc: Array<any> = [
    {
      // @ts-ignore
      name: this.$i18n.t("dash.inc"),
      data: [],
    },
    {
      // @ts-ignore
      name: this.$i18n.t("dash.realInc"),
      data: [],
    },
    {
      // @ts-ignore
      name: this.$i18n.t("dash.invoices"),
      data: [],
    },
  ];
  visHist: Array<any> = [];
  makeOrderDialog: boolean = false;
  suppliers: Supplier[] = [];
  selectedSupplier: string = "";
  role: string = "";

  mounted() {
    this.role = this.$auth.$state.user.Role;
  }

  @Watch("$route", { immediate: true, deep: true })
  onUrlChange(newVal: Route) {
    this.toggleSpinner(true);
    this.getData(this.$route.query.storeId);
  }

  getData(storeId: any) {
    if (storeId !== undefined) {
      this.ord = [];
      this.value = [];
      this.inc = [];
      this.histVisit = [];
      this.histOrd = [];
      this.histInc = [];
      this.products = [];
      this.invoices = [];
      if (!this.calculating) {
        this.calculating = true;
        this.getTrends(storeId);
      }
    } else {
      this.toggleSpinner(false);
    }
  }

  getProducts(storeId: any) {
    this.$axios
      .get("/products/needs/" + storeId + "/500/0")
      .then((response: any) => {
        if (response.data.success === true && response.data.products !== null) {
          this.products = response.data.products.filter(
            (product: ProductToStore) => product.Quantity > 0
          );
          this.toggleSpinner(false);
          this.calculating = false;
        } else {
          this.toggleSpinner(false);
        }
      });
  }

  getTrends(storeId: any) {
    this.$axios
      .get("/prediction/query/visitors/" + storeId + "/2/2")
      .then((response: any) => {
        if (response.data.success === true) {
          const histVisit = response.data.prediction
            .filter((item: any) => item.Index === "d0")
            .map((record: any) => [record.Date.split(" ")[0], record.Val]);
          let histValue: number[] = [];
          histVisit.forEach((val: any[], d: number) => {
            histValue.push(parseInt(val[1]));
            if (d % 10 === 0) {
              const sum = histValue.reduce((a, b) => a + b, 0);
              const avg = sum / histValue.length || 0;
              this.histVisit.push([val[0], avg]);
              histValue = [];
            }
          });

          let value = [];
          for (let d = 0; d < 60; d++) {
            let date = new Date();
            date.setDate(date.getDate() + d);
            const day = date.toISOString().slice(0, 10);
            const base = response.data.prediction.filter(
              (item: any) => item.Date.split(" ")[0] === day
            );
            if (base.length > 0) {
              value.push(parseInt(base[0].Val));
            } else {
              value.push(0);
            }

            if (d % 10 === 0) {
              const sum = value.reduce((a, b) => a + b, 0);
              const avg = sum / value.length || 0;
              this.value.push([date.getTime(), avg]);
              value = [];
            }
          }
          this.$refs.realtimeVisitors
            .updateSeries([{ data: this.value }, { data: this.histVisit }])
            .then(() => this.getOrders(storeId));
        } else {
          if (response.data.error !== undefined) {
            this.updateText(response.data.error);
            this.updateColor("red");
            this.updateShow(true);
          }
          this.toggleSpinner(false);
          this.waitingForData = false;
        }
      })
      .catch((error: any) => {
        console.log(error.message);
        this.toggleSpinner(false);
        this.waitingForData = false;
      });
  }

  async getOrders(storeId: any) {
    await this.getInvoices(storeId);
    this.$axios
      .get("/prediction/query/orders/" + storeId + "/2/2")
      .then((response: any) => {
        if (response.data.success === true) {
          this.histOrd = response.data.prediction
            .filter((item: any) => item.Index === "d0")
            .map((record: any) => [record.Date.split(" ")[0], record.Val]);
          this.histInc = response.data.prediction
            .filter((item: any) => item.Index === "d0")
            .map((record: any) => [
              record.Date.split(" ")[0],
              Math.ceil(record.Val * response.data.saoa),
            ]);
          for (let d = 0; d < 60; d++) {
            let date = new Date();
            date.setDate(date.getDate() + d);
            const day = date.toISOString().slice(0, 10);
            const base = response.data.prediction.filter(
              (item: any) => item.Date.split(" ")[0] === day
            );
            if (base.length > 0) {
              let min = Math.ceil(
                Math.min(...base.map((item: any) => item.Val))
              );
              this.ord.push([date.getTime(), min]);
              this.inc.push([
                date.getTime(),
                Math.ceil(min * response.data.saoa),
              ]);
            }
          }
          this.$refs.realtimeOrd
            .updateSeries([{ data: this.ord }, { data: this.histOrd }])
            .then(() => this.getProducts(storeId));
          this.$refs.realtimeInc.updateSeries([
            { data: this.inc },
            { data: this.histInc },
            {data: this.invoices}
          ]);
        } else {
          if (response.data.error !== undefined) {
            this.updateText(response.data.error);
            this.updateColor("red");
            this.updateShow(true);
          }
          this.toggleSpinner(false);
        }
      });
  }

  getInvoices(storeId: any) {
    this.$axios
      .get("/invoice/query/" + storeId + "/0/2")
      .then((response: any) => {
        if (response.data.success === true) {
          response.data.invoices.forEach((item: any) => {
            let time = new Date(item.DueDate).getTime();
            this.invoices.push([time, item.Amount]);
          });
        } else {
          if (response.data.error !== undefined) {
            this.updateText(response.data.error);
            this.updateColor("red");
            this.updateShow(true);
          }
          this.toggleSpinner(false);
        }
      });
  }

  signOut() {
    this.$auth.logout();
    this.$nuxt.$options.router.push("/sign");
  }

  solved(item: Product) {
    console.log(item);
  }

  makeOrder() {
    this.$axios
      .get("/supplier/" + this.$route.query.storeId)
      .then((response: any) => {
        if (response.data.success === true) {
          this.suppliers = response.data.suppliers;
          this.makeOrderDialog = true;
        }
      });
  }

  closeMakeOrder() {
    this.makeOrderDialog = false;
    this.selectedSupplier = "";
  }

  sendOrderToSupplier() {
    const supplier = this.suppliers.find(
      (item) => item.Id === this.selectedSupplier
    );
    let htmlTemplate = supplier?.Template + "<br>";
    this.selected.forEach((product) => {
      htmlTemplate += product.code + ": " + product.number + " Qty <br>";
    });
    this.$axios
      .post("/supplier/order", {
        recipient: supplier?.Email,
        subject: supplier?.Subject,
        template: htmlTemplate,
      })
      .then((response: any) => {
        if (response.data.success === true) {
          this.makeOrderDialog = false;
          this.selectedSupplier = "";
          this.updateText(this.$i18n.t("index.orderSend"));
          this.updateColor("green");
          this.updateShow(true);
        }
      });
  }

  filterByDate(value: any, search: string, item: ProductToStore) {
    if (this.productFilter) {
      return item.DateToNeed === search + "T00:00:00Z" && item.ProductCode === this.productFilter;
    } else {
      return item.DateToNeed === search + "T00:00:00Z";
    }
  }

  clearDate() {
    this.date = "";
    this.productFilter = "";
  }
}
</script>

<style scoped>
.v-card__title {
  color: #fff !important;
  background: #4633af !important;
}
.v-card__title.white--text {
  color: #fff !important;
  background-color: #4633af !important;
}
.redColor {
  color: red;
}
.greenColor {
  color: green;
}

.theme--dark .v-card__title {
  background-color: #4f5249 !important;
}
</style>
