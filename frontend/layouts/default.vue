<template>
  <v-app :dark="dark">
    <img src="/sp_logo.svg" class="sp d-flex d-sm-none" />
    <v-navigation-drawer
      v-model="drawer"
      :mini-variant="miniVariant"
      :clipped="clipped"
      fixed
      app
    >
      <v-list>
        <v-list-item
          v-for="(item, i) in items"
          :key="i"
          :to="localePath(item.to)"
          v-model="dark"
          :label="$t(item.title)"
          router
          exact
          :class="item.to === $route.path ? 'highlighted' : ''"
        >
          <v-list-item-action>
            <v-tooltip bottom>
              <template v-slot:activator="{ on, attrs }">
                <v-icon v-bind="attrs" v-on="on">{{ item.icon }}</v-icon>
              </template>
              <span>{{ $t(item.title) }}</span>
            </v-tooltip>
          </v-list-item-action>
          <v-list-item-content>
            <v-list-item-title v-text="$t(item.title)" />
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>
    <v-app-bar
      :clipped-left="clipped"
      fixed
      app
      elevate-on-scroll
      scroll-target="#scrolling-techniques"
    >
      <img src="/sp_logo.svg" class="sp d-none d-sm-flex" />
      <v-app-bar-nav-icon @click.stop="drawer = !drawer"></v-app-bar-nav-icon>
      <v-tooltip bottom>
        <template v-slot:activator="{ on, attrs }" class="d-none d-sm-flex">
          <v-btn
            v-bind="attrs"
            v-on="on"
            icon
            @click.stop="miniVariant = !miniVariant"
          >
            <v-icon
              >mdi-{{ `chevron-${miniVariant ? "right" : "left"}` }}</v-icon
            >
          </v-btn>
        </template>
        <span>{{ $t("tooltips.collapse") }}</span>
      </v-tooltip>
      <v-divider vertical class="d-none d-sm-flex"></v-divider>
      <v-dialog v-model="dialogCr" max-width="500px">
        <template v-slot:activator="{ on: dialogCr, attrs }">
          <v-tooltip bottom class="d-none d-md-flex">
            <template
              v-slot:activator="{ on: crHelp, attrs }"
              class="d-none d-sm-flex"
            >
              <v-col
                v-bind="attrs"
                v-on="{ ...crHelp, ...dialogCr }"
                class="inf d-none d-md-flex clickable"
                cols="1"
                >
                  <table>
                    <tr>
                      <td>
                        <span>{{ $t("info.cr") }}:</span>
                      </td>
                    </tr>
                    <tr>
                      <td>
                        {{ cr }}%
                      </td>
                    </tr>
                  </table>
              </v-col>
            </template>
            <span>{{ $t("tooltips.cr") }}</span>
          </v-tooltip>
        </template>
        <v-card>
          <v-card-title>
            <span class="text-h5">{{ $t("info.cr") }}</span>
          </v-card-title>
          <v-card-text>
            <v-container><div v-html="crHelp[$i18n.locale]"></div> </v-container>
          </v-card-text>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn color="blue darken-1" text @click="dialogCr = false">
              {{ $t('index.cancel') }}
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
      <v-divider class="d-none d-sm-flex" vertical></v-divider>
      <v-dialog v-model="dialogLead" max-width="500px">
        <template v-slot:activator="{ on: dialogLead, attrs }">
          <v-tooltip bottom class="d-none d-md-flex">
            <template
              v-slot:activator="{ on: leadHelp, attrs }"
              class="d-none d-sm-flex"
            >
              <v-col
                v-bind="attrs"
                v-on="{ ...leadHelp, ...dialogLead }"
                class="inf d-none d-md-flex clickable"
                cols="1"
                >
                <table>
                    <tr>
                      <td>
                        <span>{{ $t("info.lead") }}:</span>
                      </td>
                    </tr>
                    <tr>
                      <td>
                        {{ lead }} {{ actualStore.setting.symbol }}
                      </td>
                    </tr>
                  </table>
                </v-col>
            </template>
            <span>{{ $t("tooltips.lead") }}</span>
          </v-tooltip>
        </template>
        <v-card>
          <v-card-title>
            <span class="text-h5">{{ $t("info.lead") }}</span>
          </v-card-title>
          <v-card-text>
            <v-container> <div v-html="leadHelp[$i18n.locale]"></div> </v-container>
          </v-card-text>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn color="blue darken-1" text @click="dialogLead = false">
              {{ $t('index.cancel') }}
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
      <v-divider vertical></v-divider>
      <v-dialog v-model="dialogR2" max-width="500px">
        <template v-slot:activator="{ on: dialogR2, attrs }">
          <v-tooltip bottom class="d-none d-md-flex">
            <template
              v-slot:activator="{ on: r2Help, attrs }"
              class="d-none d-sm-flex"
            >
              <v-col
                v-bind="attrs"
                v-on="{ ...r2Help, ...dialogR2 }"
                class="inf d-none d-md-flex clickable"
                cols="1"
                >
                  <table>
                      <tr>
                        <td>
                          <span>{{ $t("info.r2") }}:</span>
                        </td>
                      </tr>
                      <tr>
                        <td>
                          {{ r2 }}%
                        </td>
                      </tr>
                    </table>
                </v-col>
            </template>
            <span>{{ $t("tooltips.r2") }}</span>
          </v-tooltip>
        </template>
        <v-card>
          <v-card-title>
            <span class="text-h5">R-Squared</span>
          </v-card-title>
          <v-card-text>
            <v-container> <div v-html="r2Help[$i18n.locale]"></div> </v-container>
          </v-card-text>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn color="blue darken-1" text @click="dialogR2 = false">
              {{ $t('index.cancel') }}
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
      <v-divider vertical></v-divider>
      <v-spacer />
      <v-menu
        bottom
        class="d-none d-sm-flex"
        origin="center center"
        transition="scale-transition"
        v-if="actualStore !== null"
      >
        <template v-slot:activator="{ on, attrs }">
          <v-btn
            color="secondary d-none d-sm-flex"
            dark
            v-bind="attrs"
            v-on="on"
          >
            {{ actualStore.name }}
          </v-btn>
        </template>

        <v-list class="d-none d-sm-block">
          <v-list-item
            v-for="(item, i) in stores"
            :key="i"
            @click="setStore(item)"
          >
            <v-list-item-title>{{ item.name }}</v-list-item-title>
          </v-list-item>
        </v-list>
      </v-menu>
      <v-menu bottom origin="center center" transition="scale-transition">
        <template v-slot:activator="{ on, attrs }">
          <v-btn
            color="secondary d-none d-sm-flex"
            dark
            v-bind="attrs"
            v-on="on"
          >
            <country-flag
              :country="$i18n.locale.replace('en', 'us')"
              size="small"
            />
            {{ $i18n.locale }}
          </v-btn>
        </template>

        <v-list>
          <v-list-item
            v-for="(item, i) in languages"
            :key="i"
            @click="onLangChange(item)"
          >
            <v-list-item-title
              ><country-flag :country="item.replace('en', 'us')" size="small" />
              {{ item }}</v-list-item-title
            >
          </v-list-item>
        </v-list>
      </v-menu>
      <v-dialog v-model="dialogExcel" persistent max-width="600px">
        <template v-slot:activator="{ on: dialog, attrs }">
          <v-tooltip bottom>
            <template v-slot:activator="{ on: tooltip }">
              <v-btn icon v-bind="attrs" v-on="{ ...tooltip, ...dialog }">
                <v-icon>mdi-file-excel</v-icon>
              </v-btn>
            </template>
            <span>{{ $t("tooltips.excel") }}</span>
          </v-tooltip>
        </template>
        <v-card>
          <v-card-title>
            <span class="text-h5">{{ $t("report.excel") }}</span>
          </v-card-title>
          <v-card-text>
            <v-container>
              <v-row>
                <v-col cols="12">
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
                        class="pt-5 mr-2"
                      ></v-text-field>
                    </template>
                    <v-date-picker
                      v-model="date"
                      @input="menu = false"
                    ></v-date-picker>
                  </v-menu>
                </v-col>
              </v-row>
            </v-container>
          </v-card-text>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn color="blue darken-1" text @click="dialogExcel = false">
              {{ $t("report.close") }}
            </v-btn>
            <v-btn color="primary" @click="export2Excel">
              {{ $t("report.send") }}
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
      <v-dialog v-model="dialogPdf" persistent max-width="600px">
        <template v-slot:activator="{ on: dialog, attrs }">
          <v-tooltip bottom>
            <template v-slot:activator="{ on: tooltip, attrs }">
              <v-btn icon v-bind="attrs" v-on="{ ...tooltip, ...dialog }">
                <v-icon>mdi-file-pdf-box</v-icon>
              </v-btn>
            </template>
            <span>{{ $t("tooltips.pdf") }}</span>
          </v-tooltip>
        </template>
        <v-card>
          <v-card-title>
            <span class="text-h5">{{ $t("report.pdf") }}</span>
          </v-card-title>
          <v-card-text>
            <v-container>
              <v-row>
                <v-col cols="12">
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
                        class="pt-5 mr-2"
                      ></v-text-field>
                    </template>
                    <v-date-picker
                      v-model="date"
                      @input="menu = false"
                    ></v-date-picker>
                  </v-menu>
                </v-col>
              </v-row>
            </v-container>
          </v-card-text>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn color="blue darken-1" text @click="dialogPdf = false">
              {{ $t("report.close") }}
            </v-btn>
            <v-btn color="primary" @click="export2Pdf">
              {{ $t("report.send") }}
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
      <v-dialog v-model="shareDialog" persistent max-width="600px">
        <template v-slot:activator="{ on: shareDialog, attrs }">
          <v-tooltip bottom>
            <template v-slot:activator="{ on: tooltip, attrs }">
              <v-btn icon v-on="{ ...tooltip, ...shareDialog }">
                <v-icon>mdi-export-variant</v-icon>
              </v-btn>
            </template>
            <span>{{ $t("tooltips.share") }}</span>
          </v-tooltip>
        </template>
        <v-card>
          <v-card-title>
            <span class="text-h5">{{ $t("report.share") }}</span>
          </v-card-title>
          <v-card-text>
            <v-container>
              <v-row>
                <v-col cols="12">
                  <v-text-field
                    v-model="reportEmail"
                    :label="$t('report.email')"
                    required
                  ></v-text-field>
                </v-col>
                <v-col cols="12">
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
                        v-model="shareDate"
                        :label="$t('index.dateToStore')"
                        prepend-icon="mdi-calendar"
                        readonly
                        v-bind="attrs"
                        v-on="on"
                        class="pt-5 mr-2"
                      ></v-text-field>
                    </template>
                    <v-date-picker
                      v-model="shareDate"
                      @input="menu = false"
                    ></v-date-picker>
                  </v-menu>
                </v-col>
                <v-col cols="12" sm="6">
                  <v-checkbox
                    v-model="checkboxExcel"
                    :label="$t('report.excel')"
                  ></v-checkbox>
                </v-col>
                <v-col cols="12" sm="6">
                  <v-checkbox
                    v-model="checkboxPdf"
                    :label="$t('report.pdf')"
                  ></v-checkbox>
                </v-col>
              </v-row>
            </v-container>
          </v-card-text>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn color="blue darken-1" text @click="shareDialog = false">
              {{ $t("report.close") }}
            </v-btn>
            <v-btn color="primary" @click="shareReport">
              {{ $t("report.send") }}
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
      <v-tooltip bottom>
        <template v-slot:activator="{ on, attrs }">
          <div v-bind="attrs" v-on="on">
            <v-switch
              class="pt-5"
              v-model="dark"
              @click="switchTemplate"
            ></v-switch>
          </div>
        </template>
        <span>{{ $t("tooltips.template") }}</span>
      </v-tooltip>
      <v-tooltip bottom>
        <template v-slot:activator="{ on, attrs }">
          <v-btn icon @click.stop="logOut()" v-bind="attrs" v-on="on">
            <v-icon>mdi-export</v-icon>
          </v-btn>
        </template>
        <span>{{ $t("tooltips.logout") }}</span>
      </v-tooltip>
    </v-app-bar>
    <v-main>
      <v-container fluid id="scrolling-techniques">
        <NuxtChild :store="actualStore" />
      </v-container>
    </v-main>
    <v-footer :absolute="fixed" app>
      <v-menu
        bottom
        class="d-flex d-sm-none"
        origin="center center"
        transition="scale-transition"
        v-if="actualStore !== null"
      >
        <template v-slot:activator="{ on, attrs }">
          <v-btn
            color="secondary d-flex d-sm-none"
            dark
            v-bind="attrs"
            v-on="on"
          >
            {{ actualStore.name }}
          </v-btn>
        </template>

        <v-list class="d-flex d-sm-none">
          <v-list-item
            v-for="(item, i) in stores"
            :key="i"
            @click="setStore(item)"
          >
            <v-list-item-title>{{ item.name }}</v-list-item-title>
          </v-list-item>
        </v-list>
      </v-menu>
      <v-menu bottom origin="center center" transition="scale-transition">
        <template v-slot:activator="{ on, attrs }">
          <v-btn
            color="secondary d-flex d-sm-none"
            dark
            v-bind="attrs"
            v-on="on"
          >
            <country-flag
              :country="$i18n.locale.replace('en', 'us')"
              size="small"
            />
            {{ $i18n.locale }}
          </v-btn>
        </template>

        <v-list>
          <v-list-item
            v-for="(item, i) in languages"
            :key="i"
            @click="onLangChange(item)"
          >
            <v-list-item-title
              ><country-flag :country="item.replace('en', 'us')" size="small" />
              {{ item }}</v-list-item-title
            >
          </v-list-item>
        </v-list>
      </v-menu>
      <span>&copy; {{ new Date().getFullYear() }} storePredictor v1.0.5</span>
    </v-footer>
    <Snackbar />
    <Spinner />
  </v-app>
</template>

<script lang="ts">
import { Component, Vue, Watch } from "nuxt-property-decorator";
import Route from "~/model/Route";
import Store from "~/model/Store";
import IResponseStore from "~/model/IResponseStore";
import StoreBO from "~/model/StoreBO";
import Snackbar from "~/components/Snackbar.vue";
import Spinner from "~/components/Spinner.vue";
import { namespace } from "vuex-class";
import IDictionary from "~/model/IDictionary";
import IResponseText from "~/model/IResponseText";
import Upgrade from "~/components/Upgrade.vue";
import { JSX } from "@babel/types";
import IPlan from "~/model/IPlan";
const snackbar = namespace("Snackbar");
const spinner = namespace("Spinner");

@Component({
  components: {
    Snackbar,
    Spinner,
    Upgrade
  },
  middleware: "auth",
})
export default class DefaultLayout extends Vue {
  @snackbar.Action
  public updateText!: (newText: string) => void;

  @snackbar.Action
  public updateColor!: (newColor: string) => void;

  @snackbar.Action
  public updateShow!: (newShow: boolean) => void;

  @spinner.Action
  public toggleSpinner!: (newShow: boolean) => void;

  checkboxExcel: boolean = true;
  checkboxPdf: boolean = true;
  reportEmail: string = "";
  clipped: boolean = true;
  shareDialog: boolean = false;
  dialogPdf: boolean = false;
  dialogExcel: boolean = false;
  dialogCr: boolean = false;
  dialogLead: boolean = false;
  dialogR2: boolean = false;
  drawer: boolean = true;
  fixed: boolean = true;
  dark: boolean = false;
  lead: number = 0;
  cr: number = 0;
  r2: number = 0;
  $i18n: any;
  $t: any;
  stores: Store[] = [];
  menu: boolean = false;
  actualStore: Store = {
    name: "",
    id: "",
    setting: { date: "", currency: "", symbol: "" },
  };
  items: Route[] = [
    {
      icon: "mdi-chart-box",
      title: "menu.analytics",
      to: "/",
    },
    {
      icon: "mdi-graph",
      title: "menu.tracking",
      to: "/tracking",
    },
    {
      icon: "mdi-store-outline",
      title: "menu.stores",
      to: "/stores",
    },
    {
      icon: "mdi-apps",
      title: "menu.products",
      to: "/products",
    },
    {
      icon: "mdi-cart",
      title: "menu.orders",
      to: "/orders",
    },
    {
      icon: "mdi-account",
      title: "menu.profile",
      to: "/profile",
    },
  ];
  miniVariant: boolean = true;
  languages: Array<string> = ["en", "cz", "de", "sk", "pl", "ro", "hu"];
  id?: string;
  $route: any;
  $router: any;
  $nuxt: any;
  $axios: any;
  $config: any;
  $vuetify: any;
  $auth: any;
  shareDate: string = "";
  date: string = "";
  crHelp: IDictionary = {};
  leadHelp: IDictionary = {};
  r2Help: IDictionary = {};
  plans: IPlan[] = [];
  localePath: any;
  
  created() {
    this.$nuxt.$on("refreshStores", ($event: string) => this.getStores($event));
  }

  beforeMount() {
    if (this.$auth.strategy.options.name === 'auth0') {
      this.$axios
        .post(this.$config.internalApi + "social", {"email": this.$auth.user.email})
          .then((response: any) => {
            if (response.data.success) {
              this.$auth.setUserToken(response.data.jwt.access_token, response.data.jwt.refresh_token);
              this.$axios
                .get(this.$config.internalApi + "account/me")
                  .then((response: any) => {
                    if (response.data.success) {
                      this.$auth.setUser(response.data.user);
                      this.guard();
                      const d = new Date();
                      d.setMonth(d.getMonth() + 1);
                      this.date = d.toISOString().substr(0, 10);
                    }
                });
            }
        });
    } else {
      this.guard();
      const d = new Date();
      d.setMonth(d.getMonth() + 1);
      this.date = d.toISOString().substr(0, 10);
    }
  }

  @Watch("$route")
  onPropertyChanged(value: string, oldValue: string) {
    if (!this.$route.query.storeId && this.actualStore !== null) {
      if (this.actualStore?.id === undefined) {
        this.$router.push({
          path: this.$route.path,
          query: { storeId: this.actualStore },
        });
      } else {
        this.$router.push({
          path: this.$route.path,
          query: { storeId: this.actualStore?.id },
        });
      }
    }
  }

  async guard() {
    if (!this.$auth.loggedIn) {
      this.$nuxt.$options.router.push("/sign");
    } else {
      this.id =
        this.$auth.$state.user.Parent !== ""
          ? this.getStores(this.$auth.$state.user.Parent)
          : this.getStores(this.$auth.$state.user.Id);
      if (this.$auth.$state.user.Parent === "") {
        this.items.push({
          icon: "mdi-account-supervisor",
          title: "menu.users",
          to: "/users",
        });
      }
      //check plan
      await this.getPlans();
    }
  }

  async logOut() {
    await this.$auth.logout();
    this.$nuxt.$options.router.push("/sign");
  }

  getStores(id: string): string {
    this.stores = [];
    return this.$axios
      .get(this.$config.internalApi + "stores/" + id)
      .then((response: IResponseStore) => {
        if (response.data.success) {
          response.data.stores.forEach((element: StoreBO) => {
            this.stores.push({
              name: element.Url,
              id: element.Id,
              setting: this.$config.setting[element.CountryCode.toLowerCase()],
            });
          });
          if (!this.$route.query.storeId) {
            if (this.stores.length > 0) {
              this.actualStore = this.stores[0];
            } else {
              this.$nuxt.$options.router.push("/stores");
            }
            if (this.actualStore?.id !== "") {
              this.$router.push({
                path: this.$route.path,
                query: { storeId: this.actualStore?.id },
              });
              this.$axios
                .get("/stats/" + this.actualStore?.id)
                .then((response: any) => {
                  if (response.data.success === true) {
                    this.lead = Math.round((response.data.lead + Number.EPSILON) * 100) / 100;
                    this.cr = Math.round((response.data.cr + Number.EPSILON) * 100) / 100;
                    this.r2 = Math.round((response.data.r2 + Number.EPSILON) * 100) / 100;
                  } else {
                    this.updateText(response.data.error);
                    this.updateColor("red");
                    this.updateShow(true);
                  }
                });
              return id;
            }
          } else {
            const actualStore = this.stores.find(
              (item) => item.id === this.$route.query.storeId
            );
            if (actualStore !== undefined) {
              this.actualStore = actualStore;
            }
              this.$axios
                .get("/stats/" + this.$route.query.storeId)
                .then((response: any) => {
                  if (response.data.success === true) {
                    this.lead = Math.round((response.data.lead + Number.EPSILON) * 100) / 100;
                    this.cr = Math.round((response.data.cr + Number.EPSILON) * 100) / 100;
                    this.r2 = Math.round((response.data.r2 + Number.EPSILON) * 100) / 100;
                  } else {
                    this.updateText(response.data.error);
                    this.updateColor("red");
                    this.updateShow(true);
                  }
                });
          }
        } else {
          this.updateText(response.data.error);
          this.updateColor("red");
          this.updateShow(true);
        }
      });
  }

  switchTemplate() {
    this.$vuetify.theme.dark = !this.$vuetify.theme.dark;
  }

  setStore(item: Store) {
    this.actualStore = item;
    this.$router.push({
      path: this.$route.path,
      query: { storeId: this.actualStore.id },
    });
  }

  onLangChange(event: string) {
    this.$i18n.locale = event;
  }

  export2Pdf() {
    this.$axios
      .get(
        "/reports/pdf/" +
          this.actualStore?.id +
          "/" +
          this.date +
          "/" +
          this.$i18n.locale
      )
      .then((response: any) => {
        if (response.data.success === true && response.data.report.length > 0) {
          const url =
            this.$config.partnerApi.replace("v1", "public") +
            response.data.report;
          window.open(url, "_blank");
        } else {
          this.updateText(response.data.error);
          this.updateColor("red");
          this.updateShow(true);
        }
      });
  }

  export2Excel() {
    this.$axios
      .get("/reports/excel/" + this.actualStore?.id + "/" + this.date)
      .then((response: any) => {
        if (response.data.success === true && response.data.report.length > 0) {
          const url =
            this.$config.partnerApi.replace("v1", "public") +
            response.data.report;
          window.open(url, "_blank");
        } else {
          this.updateText(response.data.error);
          this.updateColor("red");
          this.updateShow(true);
        }
      });
  }

  shareReport() {
    this.toggleSpinner(true);
    this.$axios
      .post(
        "/reports/share/" +
          this.actualStore?.id +
          "/" +
          this.shareDate +
          "/" +
          this.$i18n.locale,
        {
          Excel: this.checkboxExcel,
          Pdf: this.checkboxPdf,
          Email: this.reportEmail,
        }
      )
      .then((response: any) => {
        if (response.data.success === true) {
          this.updateText(this.$i18n.t("report.sended"));
          this.updateColor("green");
          this.updateShow(true);
          this.shareDialog = false;
          this.toggleSpinner(false);
        } else {
          this.updateText(response.data.error);
          this.updateColor("red");
          this.updateShow(true);
          this.shareDialog = false;
          this.toggleSpinner(false);
        }
      });
  }

  loadCmsData() {
    // crHelp
      this.$axios.get(this.$config.cmsApi + "/" + this.$config.token + "text/crHelp").then((response: IResponseText) => {
        if (response.data.success) {
          this.crHelp = JSON.parse(
            response.data.text.Value
          );
        } else {
          this.updateText(response.data.error);
          this.updateColor("red");
          this.updateShow(true);
        }
      });
    // leadHelp
      this.$axios.get(this.$config.cmsApi + "/" + this.$config.token + "text/leadHelp").then((response: IResponseText) => {
        if (response.data.success) {
          this.leadHelp = JSON.parse(
            response.data.text.Value
          );
        } else {
          this.updateText(response.data.error);
          this.updateColor("red");
          this.updateShow(true);
        }
      });
      // r2Help
      this.$axios.get(this.$config.cmsApi + "/" + this.$config.token + "text/r2Help").then((response: IResponseText) => {
        if (response.data.success) {
          this.r2Help = JSON.parse(
            response.data.text.Value
          );
        } else {
          this.updateText(response.data.error);
          this.updateColor("red");
          this.updateShow(true);
        }
      });
  }

  getPlans() {
    this.$axios
      .get(this.$config.internalApi + "plans/paid")
      .then((response: any) => {
        if (response.data.success === true) {
          this.plans = response.data.plans;
          this.plans.forEach((element: IPlan) => element.text = element.Name + ' / '+ element.Price + ' ' + this.actualStore.setting?.symbol + '' + this.$i18n.t('upgrade.mothly'));
          this.$auth.$state.user.Plan = this.plans.find((plan: IPlan) => plan.Id === this.$auth.$state.user.PlanRefer);
          this.loadCmsData();
        }
      });
  }
}
</script>
<style>
.sp {
  width: 8%;
  padding-top: 0;
  padding-left: 20px;
  filter: invert() brightness(100);
}
.v-col:before {
  content: "";
  position: absolute;
  bottom: 0;
  height: 100%;
  z-index: -1;
  right: -50%;
  width: 80%;
  transform: skew(30deg);
  background-color: #007bff;
}
.v-col:after {
  content: "";
  position: absolute;
  bottom: 0;
  height: 100%;
  z-index: -1;
  right: -10%;
  width: 2000px;
}
.v-select .v-input__slot {
  margin-top: 0;
}
.inf {
  color: #fff;
}
.inf span {
  font-size: 0.7em;
}
.theme--light.v-app-bar.v-toolbar.v-sheet {
  background-color: #4633af !important;
}
.theme--light.v-btn.v-btn--icon .v-icon {
  color: #b1ea4e !important;
}
.v-application .secondary {
  background-color: #00c9db !important;
}
.apexcharts-toolbar {
  z-index: 4 !important;
}
.apexcharts-tooltip.apexcharts-theme-light {
  color: #000;
}
.apexcharts-menu-item {
  color: #000;
}
.theme--light.v-footer {
  color: #000 !important;
}
.theme--light .v-application .primary {
  color: #4633af !important;
}
.theme--light .v-list-item__action .theme--light.v-icon {
  color: #4633af !important;
}
.theme--light .v-toolbar__title {
  color: #4633af !important;
}
.theme--light.v-card__title,
.v-dialog > .v-card > .v-card__title {
  background: #4633af !important;
  color: #fff !important;
}
.v-dialog .container {
  color: #000 !important;
}
.v-data-footer__select .v-select__selections .v-select__selection--comma {
  color: #000 !important;
}
.dateSearch {
  width: 0 !important;
  color: #000 !important;
}
.dateSearch .theme--light.v-label {
  color: #000 !important;
}
.theme--light .dateSearch.v-input input {
  color: #000 !important;
}
.highlighted {
  color: #fff !important;
  background-color: #4633af;
}
.theme--dark .highlighted {
  background-color: #4f5249;
}
.highlighted .v-list-item__action .theme--light.v-icon {
  color: #fff !important;
}
.theme--dark .v-dialog .v-card__title {
  background-color: #4f5249 !important;
}
.clickable {
  cursor: pointer;
}
@media only screen and (max-width: 600px) {
  .sp {
    width: 76%;
    padding-top: 16%;
    padding-left: 25%;
    filter: none;
  }
  .v-main {
    padding: 0 !important;
  }
}
</style>
