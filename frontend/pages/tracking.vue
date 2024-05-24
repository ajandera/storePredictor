<template>
  <v-container fluid>
    <v-row>
      <v-col cols="3" offset="5">
        <v-select
          v-if="tags.length > 0"
          :items="tags"
          :label="$t('tracking.tags')"
          v-model="currentTag"
          solo
        ></v-select>
      </v-col>
      <v-col cols="2">
        <v-menu
          ref="menuFrom"
          v-model="menuFrom"
          :close-on-content-click="false"
          :return-value.sync="dateFrom"
          transition="scale-transition"
          offset-y
        >
        <template v-slot:activator="{ on, attrs }">
          <v-text-field
            v-model="dateFrom"
            label="From"
            readonly
            v-bind="attrs"
            v-on="on"
          ></v-text-field>
        </template>
        <v-date-picker 
          v-model="dateFrom"
          :type="'month'"
          min="2021-03-01"
          :max="max"
          no-title
          scrollable
        >
        </v-date-picker>
      </v-menu>
      </v-col>
      <v-col cols="2">
        <v-menu
        ref="menuTo"
        v-model="menuTo"
        :close-on-content-click="false"
        :return-value.sync="dateTo"
        transition="scale-transition"
        offset-y
      >
        <template v-slot:activator="{ on, attrs }">
          <v-text-field
            v-model="dateTo"
            label="To"
            readonly
            v-bind="attrs"
            v-on="on"
          ></v-text-field>
        </template>
        <v-date-picker 
          v-model="dateTo"
          :type="'month'"
          min="2021-03-01"
          :max="max"
          no-title
          scrollable
        >
        </v-date-picker>
      </v-menu>
      </v-col>
    </v-row>
    <v-row justify="center" align="center" v-if="waitingForData">
      <v-col cols="12">
        <v-card class="v-card-home">
          <v-card-title>
            {{ $t("tracking.visitors") }}
            <v-spacer />
            <Help title="Real Visitors" keyString="help-realVisitors" />
          </v-card-title>
          <apexchart
            ref="realtimeVisitors"
            height="300"
            type="line"
            :options="barOptions"
            :series="seriesVisitors"
          ></apexchart>
        </v-card>
      </v-col>
    </v-row>
    <v-row justify="center" align="center" v-if="waitingForData">
      <v-col cols="12">
        <v-card class="v-card-home">
          <v-card-title>
            {{ $t("tracking.orders") }}
            <v-spacer />
            <Help title="Real Orders" keyString="help-realOrders" />
          </v-card-title>
          <apexchart
            ref="realtimeOrders"
            height="300"
            type="line"
            :options="barOptions"
            :series="seriesOrders"
          ></apexchart>
        </v-card>
      </v-col>
    </v-row>
    <v-row v-if="waitingForData">
      <v-col cols="12" xs="12" sm="12" md="6">
        <v-data-table
          :headers="orderHeaders"
          :items="orders"
          class="elevation-1"
        >
          <template v-slot:top>
            <v-toolbar flat>
              <v-toolbar-title>
                {{ $t("index.orders") }}
              </v-toolbar-title>
              <v-spacer />
              <Help title="History Orders" keyString="help-historyOrders" />
            </v-toolbar>
          </template>
        </v-data-table>
      </v-col>
      <v-col cols="12" xs="12" sm="12" md="6">
        <v-data-table
          :headers="productHeaders"
          :items="products"
          class="elevation-1"
        >
          <template v-slot:top>
            <v-toolbar flat>
              <v-toolbar-title>
                {{ $t("index.buyproducts") }}
              </v-toolbar-title>
              <v-spacer />
              <Help title="Buy products" keyString="help-buyProducts" />
            </v-toolbar>
          </template>
        </v-data-table>
      </v-col>
    </v-row>
    <v-row justify="center" align="center" v-if="waitingForData">
      <v-col cols="12">
        <v-card class="v-card-home">
          <v-card-title>
            {{ $t("index.lead") }}
            <v-spacer />
            <Help title="Value of the the lead" keyString="help-leadValue" />
          </v-card-title>
          <apexchart
            ref="realtimeLead"
            height="200"
            type="line"
            :options="barOptions"
            :series="seriesLead"
          ></apexchart>
        </v-card>
      </v-col>
    </v-row>
    <v-row justify="center" align="center" v-if="waitingForData">
      <v-col cols="12">
        <v-card class="v-card-home">
          <v-card-title>
            {{ $t("index.cr") }}
            <v-spacer />
            <Help title="Conversion rate" keyString="help-conversionRate" />
          </v-card-title>
          <apexchart
            ref="realtimeCr"
            height="200"
            type="line"
            :options="barOptions"
            :series="seriesCr"
          ></apexchart>
        </v-card>
      </v-col>
    </v-row>
    <v-row justify="center" align="center" v-if="!waitingForData">
      <v-alert border="top" color="red lighten-2" dark>
        {{ $t("validation.noData") }}
      </v-alert>
    </v-row>
    <v-row justify="center" align="center">
      <v-col cols="6">
        <v-card class="v-card-home">
          <v-card-title>
            {{ $t("tracking.customerWay") }}
            <v-spacer />
            <Help title="Way of customer" keyString="help-way-customer" />
          </v-card-title>
          <vue-mermaid-string :value="diagram" />
        </v-card>
      </v-col>
      <v-col cols="6">
        <v-data-table
          :headers="wayHeaders"
          :items="way"
          class="elevation-1"
        >
          <template v-slot:top>
            <v-toolbar flat>
              <v-toolbar-title>
                {{ $t("tracking.wayProb") }}
              </v-toolbar-title>
              <v-spacer />
              <Help title="Way customers" keyString="help-wayTable" />
            </v-toolbar>
          </template>
          </v-data-table>
      </v-col>
    </v-row>
  </v-container>
</template>
<script lang="ts">
import { Component, Prop, Vue, Watch } from "nuxt-property-decorator";
import Product from "~/model/Product";
import Store from "~/model/Store";
import { namespace } from "vuex-class";
import Help from "~/components/Help.vue";
import endent from "endent";
import tinycolor from "tinycolor2";

const snackbar = namespace("Snackbar");
const spinner = namespace("Spinner");

@Component({
  components: { Help },
})
export default class TrackingPage extends Vue {
  @snackbar.Action
  public updateText!: (newText: string) => void;

  @snackbar.Action
  public updateColor!: (newColor: string) => void;

  @snackbar.Action
  public updateShow!: (newShow: boolean) => void;

  @spinner.Action
  public toggleSpinner!: (newShow: boolean) => void;

  @Prop() readonly store!: Store;
  $t: any;
  $i18n: any;
  $axios: any;
  $route: any;
  $refs: any;
  $auth: any;
  dateFrom: string = "";
  dateTo: string = new Date().toISOString().substr(0, 7);
  max: string = new Date().toISOString().substr(0, 7);
  readonly: boolean = false;
  disabled: boolean = false;
  histVisit: Array<any> = [];
  histInc: Array<any> = [];
  histOrd: Array<any> = [];
  tags: Array<string> = [""];
  currentTag: string = "";
  products: Product[] = [];
  title: string = "Tracking";
  menuFrom: boolean = false;
  menuTo: boolean = false;
  waitingForData: boolean = true;
  orderHeaders: Array<any> = [
    {
      // @ts-ignore
      text: this.$i18n.t("index.orders1"),
      align: "start",
      sortable: false,
      value: "CreatedAt",
    },
    // @ts-ignore
    { text: this.$i18n.t("index.amount"), value: "Amount" },
  ];
  orders: Array<any> = [];
  seriesVisitors: Array<any> = [
    {
      // @ts-ignore
      name: this.$i18n.t("dash.realVisit"),
      data: [],
    },
  ];
  seriesOrders: Array<any> = [
    {
      // @ts-ignore
      name: this.$i18n.t("dash.realOrders"),
      data: [],
    },
  ];
  seriesLead: Array<any> = [
    {
      // @ts-ignore
      name: this.$i18n.t("dash.lead"),
      data: [],
    },
  ];
  seriesCr: Array<any> = [
    {
      // @ts-ignore
      name: this.$i18n.t("dash.cr"),
      data: [],
    },
  ];
  barOptions: any = {
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
      size: 3,
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
        },
      ],
    },
  };
  productHeaders: Array<any> = [
    {
      // @ts-ignore
      text: this.$i18n.t("index.product"),
      align: "start",
      sortable: false,
      value: "ProductCode",
    },
    // @ts-ignore
    { text: this.$i18n.t("index.name"), value: "Name" },
    // @ts-ignore
    { text: this.$i18n.t("index.quantity"), value: "Count" },
    // @ts-ignore
    { text: this.$i18n.t("index.amount"), value: "Avg" },
  ];
  wayHeaders: Array<any> = [
    {
      // @ts-ignore
      text: this.$i18n.t("tracking.from"),
      align: "start",
      sortable: false,
      value: "From",
    },
    // @ts-ignore
    { text: this.$i18n.t("tracking.to"), value: "To" },
    // @ts-ignore
    { text: this.$i18n.t("tracking.prob"), value: "Prob" },
  ];

  way: any = [
    {"From": "Homepage", "To": "Product Detail", "Prob": "52 %"},
    {"From": "Homepage", "To": "Leave", "Prob": "21.5 %"},
    {"From": "Homepage", "To": "Categories", "Prob": "31 %"},
    {"From": "Landing Pages", "To": "Product Detail", "Prob": "15 %"},
    {"From": "Categories", "To": "Product Detail", "Prob": "75 %"},
    {"From": "Product Detail", "To": "Cart", "Prob": "14 %"},
    {"From": "Cart", "To": "Purchase", "Prob": "5 %"},
  ];

  diagram: string = endent`
      graph TD
        Contact[Contact page]
        ProductDetail[Product Detail]
        LandingPages[Landing Pages]

        ${this.style("Cart", 2)}
        ${this.style("ProductDetail", 1)}
        ${this.style("Categories", 1)}
        ${this.style("Leave", 0)}
        ${this.style("Purchase", 3)}

        Homepage --> ProductDetail
        Homepage --> Categories
        ProductDetail --> Cart
        ProductDetail --> Leave
        Homepage --> Leave
        Contact --> Leave
        Categories --> Leave
        Categories --> ProductDetail
        Cart --> Purchase
        Cart --> Leave
        LandingPages --> ProductDetail
  `;

  @Watch("dateFrom")
  onFromChanged(value: string, oldValue: string) {
    if (oldValue !== "") {
      const from = this.monthDiff(new Date(value), new Date());
      const to = this.monthDiff(new Date(this.dateTo), new Date());
      this.toggleSpinner(true);
      this.getData(this.$route.query.storeId, from, to);
      this.menuFrom = false;
      this.dateFrom = value;
    }
  }

  @Watch("dateTo")
  ontoChanged(value: string, oldValue: string) {
    if (oldValue !== "") {
      const from = this.monthDiff(new Date(this.dateFrom), new Date());
      const to = this.monthDiff(new Date(value), new Date());
      this.toggleSpinner(true);
      this.getData(this.$route.query.storeId, from, to);
      this.menuTo = false;
      this.dateTo = value;
    }
  }

  @Watch("currentTag")
  onCurrentTagChanged(value: string, oldValue: string) {
    console.log(value);
    this.toggleSpinner(true);
    this.getData(this.$route.query.storeId);
    this.menuTo = false;
  }

  mounted() {
    const date  = new Date();
    date.setMonth(date.getMonth()-2);
    this.dateFrom = date.toISOString().substr(0, 7);
    if (this.$route.query.storeId !== undefined) {
      this.toggleSpinner(true);
      this.getData(this.$route.query.storeId);
    }
  }

  @Watch("store")
  onPropertyChanged(value: Store, oldValue: Store) {
    this.toggleSpinner(true);
    this.getData(value.id);
  }

  style(node: string, colorIndex: number) {
    const colorDefs = [
      { color: "#FFE9F0", saturateStroke: -5 },
      { color: "#F3FFE9", darkenStroke: 60, saturateStroke: -5 },
      { color: "#FFF9E9", darkenStroke: 50, saturateStroke: 10 },
      { color: "#008000", darkenStroke: 40, saturateStroke: 10 }
    ];
    const colorDef = colorDefs[colorIndex];
    const fill = colorDef.color;
    const stroke = tinycolor(colorDef.color)
      .darken(colorDef.darkenStroke)
      .saturate(colorDef.saturateStroke)
      .toString();
    return `style ${node} fill: ${fill}, stroke: ${stroke}`;
  };

  getData(storeId: any, from: number = 3, to: number = 0) {
    this.histVisit = [];
    this.histInc = [];
    this.getOrdersList(storeId);
    this.getProducts(storeId);
    this.getVisitors(storeId, from, to);
  }

  head() {
    return {
      title: this.title,
    };
  }

  getOrdersList(storeId: any) {
    this.$axios
      .get(this.$config.internalApi + "orders/" + storeId + "/10/0")
      .then((response: any) => {
        if (response.data.success === true) {
          if (response.data.orders !== null) {
            this.orders = response.data.orders;
            this.orders.forEach((item) => {
              let d = new Date(item.CreatedAt);
              item.CreatedAt = d.toLocaleDateString(this.store.setting?.date);
              item.Amount = item.Amount + " " + item.Currency;
            });
          }
        } else {
          this.toggleSpinner(false);
        }
      });
  }

  getProducts(storeId: any) {
    this.$axios.get("/products/" + storeId + "/10/0").then((response: any) => {
      if (response.data.success === true) {
        if (response.data.products !== null) {
          this.products = [];
          response.data.products.forEach((item: Product) => {
            let d = new Date();
            d.setDate(d.getDate() + 4);
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
        this.toggleSpinner(false);
      }
    });
  }

  getLeadValue() {
    let lead = [];
    for (let i = 0; i < this.histVisit.length; i++) {
      const orders = this.histInc.find(
        (item) => item[0] === this.histVisit[i][0]
      );
      if (orders !== undefined) {
        lead.push([orders[0], Math.ceil(orders[1] / this.histVisit[i][1])]);
      }
    }
    this.$refs.realtimeLead.updateSeries([{ data: lead }]);
  }

  getConvRate() {
    let cr = [];
    for (let i = 0; i < this.histOrd.length; i++) {
      if (this.histVisit[i] !== undefined) {
        cr.push([
          this.histOrd[i][0],
          Math.round((this.histOrd[i][1] / this.histVisit[i][1]) * 100),
        ]);
      }
    }
    this.$refs.realtimeCr.updateSeries([{ data: cr }]);
  }

  getVisitors(storeId: any, from: number = 3, to: number = 0) {
    let dateFrom = new Date();
    // calculate from
    dateFrom.setMonth(dateFrom.getMonth() - 3);
    let dateTo = new Date();
    // calculate to
    dateTo.setMonth(dateTo.getMonth() - to);
    this.$axios
      .get("/prediction/query/visitors/" + storeId + "/" + dateFrom.toISOString().split('T')[0] + "/" + dateTo.toISOString().split('T')[0])
      .then((response: any) => {
        if (response.data.success === true) {
          this.histVisit = response.data.visitors
            .filter((item: any) => item.Tag === this.currentTag)
            .map((record: any) => [record.Day.split("T")[0], record.Visitors]);
          this.tags = response.data.visitors
            .filter((item: any) => item.Tag !== "")
            .map((record: any) => record.Tag).filter((x: any, i: any, a: any) => a.indexOf(x) === i);
          this.$refs.realtimeVisitors
            .updateSeries([{ data: this.histVisit }])
            .then(() => this.getOrders(storeId, from, to));
        } else {
          this.toggleSpinner(false);
          this.waitingForData = false;
        }
      })
      .catch((error: any) => {
        this.toggleSpinner(false);
        this.waitingForData = false;
      });
  }

  getOrders(storeId: any, from: number = 3, to: number = 0) {
    this.tags.push("");
    this.$axios
      .get("/prediction/query/orders/" + storeId + "/" + from + "/"+ to)
      .then((response: any) => {
        if (response.data.success === true) {
          const done = new Promise((resolve, reject) => {
            this.histOrd = response.data.prediction
              .filter((item: any) => item.Index === "d0")
              .map((record: any) => [record.Date.split(" ")[0], record.Val]);
            this.histInc = response.data.prediction
              .filter((item: any) => item.Index === "d0")
              .map((record: any) => [
                record.Date.split(" ")[0],
                Math.ceil(record.Val * response.data.saoa),
              ]);
            resolve(true);
          });
          done.then(() => {
            this.getConvRate();
            this.getLeadValue();
            this.$refs.realtimeOrders.updateSeries([{ data: this.histOrd }]);
            this.toggleSpinner(false);
          });
        } else {
          this.updateText(response.data.error);
          this.updateColor("red");
          this.updateShow(true);
        }
      })
      .catch((error: any) => {
        this.toggleSpinner(false);
        this.waitingForData = false;
      });
  }

  monthDiff(d1: Date, d2: Date) {
    var months;
    months = (d2.getFullYear() - d1.getFullYear()) * 12;
    months -= d1.getMonth();
    months += d2.getMonth();
    return months <= 0 ? 0 : months;
  }
}
</script>

<style scoped>
.v-icon {
  color: #000059;
}
.theme--light .v-card__title,
.v-dialog > .v-card > .v-card__title {
  background: #4633af !important;
  color: #fff !important;
}
.v-dialog .container {
  color: #000 !important;
}
</style>
