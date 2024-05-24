<template>
  <v-container fluid class="stores">
    <v-data-table
      :headers="headers"
      :items="stores"
      :items-per-page="10"
      class="elevation-1"
    >
      <template v-slot:top>
        <v-toolbar flat>
          <v-toolbar-title>{{ $t("stores.mystores") }}</v-toolbar-title>
          <v-divider class="mx-4" inset vertical></v-divider>
          <v-spacer></v-spacer>
          <v-dialog v-if="role !== 'Viewer'" v-model="dialog" max-width="500px">
            <template v-slot:activator="{ on, attrs }">
              <Help title="Store management" keyString="help-storeManagement" />
              <v-btn
                color="primary"
                dark
                class="mb-2 mr-2"
                v-bind="attrs"
                v-on="on"
              >
                {{ $t("stores.add") }}
              </v-btn>
            </template>
            <v-card>
              <v-form ref="form" v-model="valid" lazy-validation>
                <v-card-title>
                  <span class="text-h5">{{ formTitle }}</span>
                </v-card-title>
                <v-card-text>
                  <v-container>
                    <v-row>
                      <v-col cols="12" sm="6" md="6">
                        <v-text-field
                          v-model="editedItem.url"
                          :counter="50"
                          :rules="urlRules"
                          :label="$t('stores.form.url')"
                          required
                        ></v-text-field>
                      </v-col>
                      <v-col cols="12" sm="6" md="6">
                        <v-autocomplete
                          :items="countries"
                          :filter="customFilter"
                          item-text="Name"
                          item-value="Code"
                          :rules="[(v) => !!v || 'Item is required']"
                          :label="$t('stores.form.country')"
                          v-model="editedItem.country"
                        ></v-autocomplete>
                      </v-col>
                    </v-row>
                    <v-row>
                      <v-col cols="12">
                        <p>{{ $t('stores.feedInfo') }}</p>
                        <v-text-field
                          v-model="editedItem.feed"
                          :counter="100"
                          :label="$t('stores.form.feed')"
                          required
                        ></v-text-field>
                      </v-col>
                    </v-row>
                    <v-row>
                      <v-col cols="12">
                        <p>{{ $t('stores.window') }}</p>
                        <v-text-field
                          v-model="editedItem.window"
                          :counter="5"
                          :label="$t('stores.window')"
                          required
                        ></v-text-field>
                      </v-col>
                    </v-row>
                  </v-container>
                </v-card-text>
                <v-card-actions>
                  <v-spacer></v-spacer>
                  <v-btn color="blue darken-1" text @click="close">
                    Cancel
                  </v-btn>
                  <v-btn
                    color="primary"
                    dark
                    :disabled="!valid"
                    @click="addStore"
                  >
                    Save
                  </v-btn>
                </v-card-actions>
              </v-form>
            </v-card>
          </v-dialog>
          <v-dialog v-model="dialogDelete" max-width="500px">
            <v-card>
              <v-card-title class="text-h5">{{
                $t("stores.delete")
              }}</v-card-title>
              <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn color="blue darken-1" text @click="closeDelete">{{
                  $t("stores.cancel")
                }}</v-btn>
                <v-btn color="primary" @click="deleteItemConfirm">{{
                  $t("stores.ok")
                }}</v-btn>
                <v-spacer></v-spacer>
              </v-card-actions>
            </v-card>
          </v-dialog>
          <v-dialog v-model="dialogDetail">
            <v-card>
              <v-card-title class="text-h5">{{
                $t("stores.tracking")
              }}</v-card-title>
              <v-card-text>
                <v-container>
                  <v-row>
                    <v-col cols="12">
                      <h2>{{ $t("stores.website") }}</h2>
                      <br />
                      <pre>
  &lt;script&gt;
  (function(i,s,o,g,r,a,m){i['StorePredictorObjectV2']=r;i[r]=i[r]||function(){
  (i[r].q=i[r].q||[]).push(arguments)},i[r].l=1*new Date();a=s.createElement(o),
  m=s.getElementsByTagName(o)[0];a.async=1;a.src=g;m.parentNode.insertBefore(a,m)
  })(window,document,'script','https://cdn.storepredictor.com/sp.v2.js','sp');

  sp('create', '{{ spCode }}', 'auto');
  &lt;/script&gt;
                      </pre>
                      <br />
                      <h2>{{ $t("stores.connect.row1") }}</h2>
                      <p>{{ $t("stores.connect.row2") }}</p>
                      <pre>
  &lt;script&gt;sp('send', 'pageview', '');&lt;/script&gt;
                      </pre>
                      <br />
                      <p>{{ $t("stores.connect.row3") }}</p>
                      <pre>
  &lt;script&gt;sp('send', 'pageview', 'PRODUCT_CODE');&lt;/script&gt;
                      </pre>
                      <br />
                      <h2>{{ $t("stores.connect.row4") }}</h2>
                      <br />
                      <pre>
  &lt;script&gt;
  sp('create', '{{ spCode }}', 'auto');
  sp('send', 'order', { orderId: 'ORDER_NUMBER', items:   [ { productCode: 'PRODUCT_CODE', unitPrice: 100, quantity: 1 } ], totalPrice: 100, currency: 'EUR' });
  &lt;/script&gt;
                      </pre>
                    </v-col>
                  </v-row>
                  <v-row>
                    <v-col cols="12">
                      <h2>{{ $t('stores.plugins.overview') }}</h2>
                    </v-col>
                  </v-row>
                  <v-row align="center" justify="center"> 
                    <v-col cols="2">
                      <v-img src="https://storepredictor.com/integration/shopify.png" class=""></v-img>
                    </v-col>
                    <v-col cols="2">
                      <v-img src="https://storepredictor.com/integration/wordpress.png" class="gray"></v-img>
                    </v-col>
                    <v-col cols="2">
                      <v-img src="https://storepredictor.com/integration/prestashop.png"></v-img>
                    </v-col>
                    <v-col cols="2">
                      <v-img src="https://storepredictor.com/integration/magento.png"></v-img>
                    </v-col>
                    <v-col cols="2">
                      <v-img src="https://storepredictor.com/integration/opencart.png"></v-img>
                    </v-col>
                    <v-col cols="2">
                      <v-img src="https://storepredictor.com/integration/gtm.png"></v-img>
                    </v-col>
                  </v-row>
                  <v-row align="center">
                    <v-col cols="6">
                      <v-btn href="https://storepredictor.com/integrations" target="_blank">{{ $t('stores.plugins.integration') }}</v-btn>
                    </v-col>
                  </v-row>
                </v-container>
              </v-card-text>
              <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn color="blue darken-1" text @click="closeDetail">{{
                  $t("stores.close")
                }}</v-btn>
                <v-spacer></v-spacer>
              </v-card-actions>
            </v-card>
          </v-dialog>
          <v-dialog v-model="dialogSuppliers" max-width="1000px">
            <v-card>
              <v-card-title class="text-h5">{{ storeName }}</v-card-title>
              <v-card-text>
                <v-container>
                  <v-data-table
                    :headers="headersForStores"
                    :items="suppliers"
                    :items-per-page="10"
                    class="elevation-1"
                  >
                    <template v-slot:top>
                      <v-toolbar flat>
                        <v-toolbar-title>{{
                          $t("stores.suppliers")
                        }}</v-toolbar-title>
                        <v-divider class="mx-4" inset vertical></v-divider>
                        <v-spacer></v-spacer>
                        <v-dialog v-model="dialogAddSupplier" max-width="500px">
                          <template v-slot:activator="{ on, attrs }">
                            <v-btn
                              color="primary"
                              dark
                              class="mb-2"
                              v-bind="attrs"
                              v-on="on"
                            >
                              {{ $t("suppliers.add") }}
                            </v-btn>
                          </template>
                          <v-card>
                            <v-form
                              ref="supplierForm"
                              v-model="validSupplier"
                              lazy-validation
                            >
                              <v-card-title>
                                <span class="text-h5">{{
                                  formTitleSupplier
                                }}</span>
                              </v-card-title>
                              <v-card-text>
                                <v-container>
                                  <v-row>
                                    <v-col cols="6">
                                      <v-text-field
                                        v-model="supplier.Name"
                                        :counter="30"
                                        :rules="ruleName"
                                        :label="$t('suppliers.form.name')"
                                        required
                                      ></v-text-field>
                                      <v-text-field
                                        v-model="supplier.Person"
                                        :counter="30"
                                        :rules="rulePerson"
                                        :label="$t('suppliers.form.person')"
                                        required
                                      ></v-text-field>
                                      <v-text-field
                                        v-model="supplier.Email"
                                        :rules="ruleEmail"
                                        :label="$t('suppliers.form.email')"
                                        required
                                      ></v-text-field>
                                      <v-text-field
                                        v-model="supplier.Phone"
                                        :counter="10"
                                        :rules="rulePhone"
                                        :label="$t('suppliers.form.phone')"
                                        required
                                      ></v-text-field>
                                    </v-col>
                                    <v-col cols="6">
                                      <v-text-field
                                        v-model="supplier.Street"
                                        :counter="20"
                                        :rules="ruleStreet"
                                        :label="$t('suppliers.form.street')"
                                        required
                                      ></v-text-field>
                                      <v-text-field
                                        v-model="supplier.City"
                                        :counter="20"
                                        :rules="ruleCity"
                                        :label="$t('suppliers.form.city')"
                                        required
                                      ></v-text-field>
                                      <v-text-field
                                        v-model="supplier.Zip"
                                        :counter="5"
                                        :rules="ruleZip"
                                        :label="$t('suppliers.form.zip')"
                                        required
                                      ></v-text-field>
                                      <v-autocomplete
                                        :items="countries"
                                        :filter="customFilter"
                                        item-text="Name"
                                        item-value="Code"
                                        :rules="ruleCountry"
                                        :label="$t('suppliers.form.country')"
                                        v-model="supplier.Country"
                                      ></v-autocomplete>
                                    </v-col>
                                  </v-row>
                                  <v-row>
                                    <v-col cols="12">
                                      <h3>
                                        {{ $t("suppliers.form.template") }}
                                      </h3>
                                      <v-text-field
                                        v-model="supplier.Subject"
                                        :counter="50"
                                        :rules="ruleSubject"
                                        :label="$t('suppliers.form.subject')"
                                        required
                                      ></v-text-field>
                                      <quill-editor
                                        :ref="supplier.Template"
                                        v-model="supplier.Template"
                                        :options="editorOption"
                                        @blur="onEditorBlur($event)"
                                        @focus="onEditorFocus($event)"
                                        @ready="onEditorReady($event)"
                                      />
                                    </v-col>
                                  </v-row>
                                </v-container>
                              </v-card-text>
                              <v-card-actions>
                                <v-spacer></v-spacer>
                                <v-btn
                                  color="blue darken-1"
                                  text
                                  @click="closeAddSuppliers"
                                >
                                  Cancel
                                </v-btn>
                                <v-btn
                                  color="primary"
                                  dark
                                  :disabled="!validSupplier"
                                  @click="addSupplier"
                                >
                                  Save
                                </v-btn>
                              </v-card-actions>
                            </v-form>
                          </v-card>
                        </v-dialog>
                        <v-dialog
                          v-model="dialogSupplierDelete"
                          max-width="500px"
                        >
                          <v-card>
                            <v-card-title class="text-h5">{{
                              $t("suppliers.confirm")
                            }}</v-card-title>
                            <v-card-actions>
                              <v-spacer></v-spacer>
                              <v-btn
                                color="blue darken-1"
                                text
                                @click="closeSupplierDelete"
                                >{{ $t("stores.cancel") }}</v-btn
                              >
                              <v-btn
                                color="primary"
                                @click="deleteSupplierConfirm"
                                >{{ $t("stores.ok") }}</v-btn
                              >
                              <v-spacer></v-spacer>
                            </v-card-actions>
                          </v-card>
                        </v-dialog>
                      </v-toolbar>
                    </template>
                    <template v-slot:item.actions="{ item }">
                      <v-tooltip bottom>
                        <template v-slot:activator="{ on, attrs }">
                          <v-icon
                            small
                            class="mr-2"
                            v-bind="attrs"
                            v-on="on"
                            @click="editSupplier(item)"
                          >
                            mdi-pencil
                          </v-icon>
                        </template>
                        <span>{{ $t("tooltips.edit") }}</span>
                      </v-tooltip>
                      <v-tooltip bottom>
                        <template v-slot:activator="{ on, attrs }">
                          <v-icon
                            small
                            v-bind="attrs"
                            v-on="on"
                            @click="deleteSupplier(item)"
                          >
                            mdi-delete
                          </v-icon>
                        </template>
                        <span>{{ $t("tooltips.delete") }}</span>
                      </v-tooltip>
                    </template>
                    <template v-slot:no-data> </template>
                  </v-data-table>
                </v-container>
              </v-card-text>
              <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn color="blue darken-1" text @click="closeSuppliers">{{
                  $t("stores.close")
                }}</v-btn>
                <v-spacer></v-spacer>
              </v-card-actions>
            </v-card>
          </v-dialog>
          <v-dialog v-model="dialogInvoices" max-width="1000px">
            <v-card>
              <v-card-title class="text-h5">{{ storeName }}</v-card-title>
              <v-card-text>
                <v-container>
                  <v-data-table
                    :headers="headersForInvoices"
                    :items="invoices"
                    :items-per-page="10"
                    class="elevation-1"
                  >
                    <template v-slot:top>
                      <v-toolbar flat>
                        <v-toolbar-title>{{
                          $t("stores.invoices")
                        }}</v-toolbar-title>
                        <v-divider class="mx-4" inset vertical></v-divider>
                        <v-spacer></v-spacer>
                        <v-dialog v-model="dialogAddInvoice" max-width="500px">
                          <template v-slot:activator="{ on: invoice, attrs }">
                            <v-btn
                              color="primary"
                              dark
                              class="mb-2"
                              v-bind="attrs"
                              v-on="{ ...invoice }"
                            >
                              {{ $t("invoices.add") }}
                            </v-btn>
                          </template>
                          <v-card>
                            <v-form
                              ref="invoiceForm"
                              v-model="validInvoice"
                              lazy-validation
                            >
                              <v-card-title>
                                <span class="text-h5">{{
                                  formTitleInvoice
                                }}</span>
                              </v-card-title>
                              <v-card-text>
                                <v-container>
                                  <v-row>
                                    <v-col cols="12">
                                      <v-menu
                                        v-model="menu"
                                        :close-on-content-click="true"
                                        :nudge-right="40"
                                        transition="scale-transition"
                                        offset-y
                                        min-width="auto"
                                      >
                                        <template
                                          v-slot:activator="{ on, attrs }"
                                        >
                                          <v-text-field
                                            v-model="invoice.DueDate"
                                            :label="$t('invoices.form.duedate')"
                                            prepend-icon="mdi-calendar"
                                            readonly
                                            v-bind="attrs"
                                            v-on="on"
                                            class="pt-5 mr-2"
                                            :rules="ruleDate"
                                          ></v-text-field>
                                        </template>
                                        <v-date-picker
                                          v-model="invoice.DueDate"
                                          @input="menu = false"
                                        ></v-date-picker>
                                      </v-menu>
                                      <v-text-field
                                        v-model="invoice.Amount"
                                        :counter="10"
                                        :rules="ruleAmount"
                                        :label="$t('invoices.form.amount')"
                                        required
                                      ></v-text-field>
                                      <v-text-field
                                        v-model="invoice.Currency"
                                        :counter="3"
                                        :rules="ruleCurrency"
                                        :label="$t('invoices.form.currency')"
                                        required
                                      ></v-text-field>
                                    </v-col>
                                  </v-row>
                                </v-container>
                              </v-card-text>
                              <v-card-actions>
                                <v-spacer></v-spacer>
                                <v-btn
                                  color="blue darken-1"
                                  text
                                  @click="closeAddInvoice"
                                >
                                  Cancel
                                </v-btn>
                                <v-btn
                                  color="primary"
                                  dark
                                  :disabled="!validInvoice"
                                  @click="addInvoice"
                                >
                                  Save
                                </v-btn>
                              </v-card-actions>
                            </v-form>
                          </v-card>
                        </v-dialog>
                      </v-toolbar>
                    </template>
                    <template v-slot:item.DueDate="{ item }">
                      <span>{{
                        new Date(item.DueDate).toLocaleDateString()
                      }}</span>
                    </template>
                    <template v-slot:item.actions="{ item }">
                      <v-tooltip bottom>
                        <template v-slot:activator="{ on, attrs }">
                          <v-icon
                            small
                            class="mr-2"
                            v-bind="attrs"
                            v-on="on"
                            @click="editInvoice(item)"
                          >
                            mdi-pencil
                          </v-icon>
                        </template>
                        <span>{{ $t("tooltips.edit") }}</span>
                      </v-tooltip>
                      <v-tooltip bottom>
                        <template v-slot:activator="{ on, attrs }">
                          <v-icon
                            small
                            v-bind="attrs"
                            v-on="on"
                            @click="deleteInvoice(item)"
                          >
                            mdi-delete
                          </v-icon>
                        </template>
                        <span>{{ $t("tooltips.delete") }}</span>
                      </v-tooltip>
                    </template>
                    <template v-slot:no-data> </template>
                  </v-data-table>
                </v-container>
              </v-card-text>
              <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn color="blue darken-1" text @click="closeInvoices">{{
                  $t("stores.close")
                }}</v-btn>
                <v-spacer></v-spacer>
              </v-card-actions>
            </v-card>
          </v-dialog>
          <v-dialog
            v-model="dialogInvoiceDelete"
            max-width="500px"
          >
            <v-card>
              <v-card-title class="text-h5">{{
                $t("invoices.confirm")
              }}</v-card-title>
              <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn
                  color="blue darken-1"
                  text
                  @click="closeInvoiceDelete"
                  >{{ $t("stores.cancel") }}</v-btn
                >
                <v-btn
                  color="primary"
                  @click="deleteInvoiceConfirm"
                  >{{ $t("stores.ok") }}</v-btn
                >
                <v-spacer></v-spacer>
              </v-card-actions>
            </v-card>
          </v-dialog>
        </v-toolbar>
      </template>
      <template v-slot:item.actions="{ item }">
        <v-btn-toggle rounded dense>
          <v-tooltip bottom>
            <template v-slot:activator="{ on, attrs }">
              <v-btn v-bind="attrs" v-on="on" @click="detailItem(item)">
                <v-icon small class="mr-2" color="success"> mdi-web </v-icon>
              </v-btn>
            </template>
            <span>{{ $t("tooltips.implement") }}</span>
          </v-tooltip>
          <v-tooltip bottom v-if="role !== 'Viewer'">
            <template v-slot:activator="{ on, attrs }">
              <v-btn v-bind="attrs" v-on="on" @click="showSuppliers(item)">
                <v-icon small> mdi-factory </v-icon>
              </v-btn>
            </template>
            <span>{{ $t("tooltips.suppliers") }}</span>
          </v-tooltip>
          <v-tooltip bottom v-if="role !== 'Viewer'">
            <template v-slot:activator="{ on, attrs }">
              <v-btn v-bind="attrs" v-on="on" @click="showInvoices(item)">
                <v-icon small> mdi-file </v-icon>
              </v-btn>
            </template>
            <span>{{ $t("tooltips.invoices") }}</span>
          </v-tooltip>
          <v-tooltip bottom v-if="role !== 'Viewer'">
            <template v-slot:activator="{ on, attrs }">
              <v-btn v-bind="attrs" v-on="on" @click="editItem(item)">
                <v-icon small class="mr-2" color="secondary">
                  mdi-pencil
                </v-icon>
              </v-btn>
            </template>
            <span>{{ $t("tooltips.edit") }}</span>
          </v-tooltip>
          <v-tooltip bottom v-if="role !== 'Viewer'">
            <template v-slot:activator="{ on, attrs }">
              <v-btn v-bind="attrs" v-on="on" @click="deleteItem(item)">
                <v-icon small color="error"> mdi-delete </v-icon>
              </v-btn>
            </template>
            <span>{{ $t("tooltips.delete") }}</span>
          </v-tooltip>
        </v-btn-toggle>
      </template>
      <template v-slot:no-data> </template>
    </v-data-table>
  </v-container>
</template>

<script lang="ts">
import { Component, Vue, Watch } from "nuxt-property-decorator";
import Invoice from "~/model/Invoice";
import Store from "~/model/Store";
import StoreBO from "~/model/StoreBO";
import Supplier from "~/model/Supplier";
import { namespace } from "vuex-class";
import Help from "~/components/Help.vue";
import ICountry from "~/model/ICountry";

const snackbar = namespace("Snackbar");
const spinner = namespace("Spinner");

@Component({
  components: { Help },
})
export default class StoresPage extends Vue {
  @snackbar.Action
  public updateText!: (newText: string) => void;

  @snackbar.Action
  public updateColor!: (newColor: string) => void;

  @snackbar.Action
  public updateShow!: (newShow: boolean) => void;

  @spinner.Action
  public toggleSpinner!: (newShow: boolean) => void;

  $i18n: any;
  stores: StoreBO[] = [];
  dialog: boolean = false;
  dialogDelete: boolean = false;
  dialogDetail: boolean = false;
  dialogSuppliers: boolean = false;
  dialogInvoices: boolean = false;
  spCode: string = "";
  storeName: string | undefined = "";
  headers: Array<any> = [
    {
      // @ts-ignore
      text: this.$i18n.t("stores.form.url"),
      align: "start",
      sortable: false,
      value: "Url",
    },
    // @ts-ignore
    { text: this.$i18n.t("stores.country"), value: "CountryCode" },
    // @ts-ignore
    { text: this.$i18n.t("stores.feed"), value: "XmlFeed" },
    // @ts-ignore
    { text: this.$i18n.t("stores.actions"), value: "actions", sortable: false },
  ];
  editedIndex: number = -1;
  editedItem: Store = {
    id: "",
    url: "",
    country: "SK",
    feed: "",
    window: 90
  };
  headersForStores: Array<any> = [
    {
      // @ts-ignore
      text: this.$i18n.t("suppliers.form.name"),
      align: "start",
      sortable: false,
      value: "Name",
    },
    // @ts-ignore
    { text: this.$i18n.t("suppliers.form.street"), value: "Street" },
    // @ts-ignore
    { text: this.$i18n.t("suppliers.form.city"), value: "City" },
    // @ts-ignore
    { text: this.$i18n.t("suppliers.form.email"), value: "Email" },
    // @ts-ignore
    { text: this.$i18n.t("suppliers.form.phone"), value: "Phone" },
    // @ts-ignore
    { text: this.$i18n.t("suppliers.form.person"), value: "Person" },
    {
      // @ts-ignore
      text: this.$i18n.t("suppliers.form.actions"),
      value: "actions",
      sortable: false,
    },
  ];
  headersForInvoices: Array<any> = [
    {
      // @ts-ignore
      text: this.$i18n.t("invoices.form.duedate"),
      align: "start",
      sortable: false,
      value: "DueDate",
    },
    // @ts-ignore
    { text: this.$i18n.t("invoices.form.amount"), value: "Amount" },
    // @ts-ignore
    { text: this.$i18n.t("invoices.form.currency"), value: "Currency" },
    {
      // @ts-ignore
      text: this.$i18n.t("invoices.form.actions"),
      value: "actions",
      sortable: false,
    },
  ];
  // @ts-ignore
  title: string = this.$i18n.t("stores.title");
  valid: boolean = true;
  urlRules: Array<any> = [
    // @ts-ignore
    (v: any) => !!v || this.$i18n.t("validation.url"),
    // @ts-ignore
    (v: string | any[]) =>
      (v && v.length <= 50) || this.$i18n.t("validation.url"),
  ];
  countries: Array<ICountry> = [];
  $axios: any;
  $t: any;
  suppliers: Supplier[] = [];
  supplier: Supplier = {
    Id: "",
    Name: "",
    Street: "",
    City: "",
    Zip: "",
    Email: "",
    Phone: "",
    Person: "",
    Country: "",
    Template: "",
    Subject: "",
  };
  invoices: Invoice[] = [];
  invoice: Invoice = {
    Id: "",
    DueDate: "",
    Amount: 0,
    Currency: "",
  };
  supplierStore: StoreBO = this.stores[0];
  dialogSupplierDelete: boolean = false;
  dialogAddSupplier: boolean = false;
  validSupplier: boolean = false;
  dialogInvoiceDelete: boolean = false;
  template: string = "";
  editorOption: any = {
    // Some Quill options...
  };
  ruleName: Array<any> = [
    // @ts-ignore
    (v: any) => !!v || this.$i18n.t("validation.name"),
    // @ts-ignore
    (v: string | any[]) =>
      (v && v.length <= 30) || this.$i18n.t("validation.name"),
  ];
  rulePerson: Array<any> = [
    // @ts-ignore
    (v: any) => !!v || this.$i18n.t("validation.person"),
    // @ts-ignore
    (v: string | any[]) =>
      (v && v.length <= 30) || this.$i18n.t("validation.person"),
  ];
  ruleStreet: Array<any> = [
    // @ts-ignore
    (v: any) => !!v || this.$i18n.t("validation.street"),
    // @ts-ignore
    (v: string | any[]) =>
      (v && v.length <= 20) || this.$i18n.t("validation.street"),
  ];
  ruleCity: Array<any> = [
    // @ts-ignore
    (v: any) => !!v || this.$i18n.t("validation.city"),
    // @ts-ignore
    (v: string | any[]) =>
      (v && v.length <= 20) || this.$i18n.t("validation.city"),
  ];
  ruleEmail: Array<any> = [
    (v: any) => !!v || this.$i18n.t("validation.email"),
    (v: string) => /.+@.+\..+/.test(v) || this.$i18n.t("validation.email"),
  ];
  ruleZip: Array<any> = [
    // @ts-ignore
    (v: any) => !!v || this.$i18n.t("validation.zip"),
    // @ts-ignore
    (v: string | any[]) =>
      (v && v.length <= 5) || this.$i18n.t("validation.zip"),
  ];
  rulePhone: Array<any> = [
    // @ts-ignore
    (v: any) => !!v || this.$i18n.t("validation.phone"),
    // @ts-ignore
    (v: string | any[]) =>
      (v && v.length <= 10) || this.$i18n.t("validation.phone"),
  ];
  ruleCountry: Array<any> = [
    // @ts-ignore
    (v: any) => !!v || this.$i18n.t("validation.country"),
    // @ts-ignore
    (v: string | any[]) =>
      (v && v.length <= 15) || this.$i18n.t("validation.country"),
  ];
  ruleSubject: Array<any> = [
    // @ts-ignore
    (v: any) => !!v || this.$i18n.t("validation.subject"),
    // @ts-ignore
    (v: string | any[]) =>
      (v && v.length <= 20) || this.$i18n.t("validation.subject"),
  ];
  dialogAddInvoice: boolean = false;
  validInvoice: boolean = false;
  menu: boolean = false;
  ruleDate: Array<any> = [
    // @ts-ignore
    (v: any) => !!v || this.$i18n.t("validation.date"),
  ];
  ruleAmount: Array<any> = [
    // @ts-ignore
    (v: any) => !!v || this.$i18n.t("validation.amount"),
    // @ts-ignore
    (v: string | any[]) =>
      (v && v.length <= 10) || this.$i18n.t("validation.amount"),
  ];
  ruleCurrency: Array<any> = [
    // @ts-ignore
    (v: any) => !!v || this.$i18n.t("validation.currency"),
    // @ts-ignore
    (v: string | any[]) =>
      (v && v.length <= 3) || this.$i18n.t("validation.currency"),
  ];
  $refs: any;
  role: string = "";
  $content: any;

  head() {
    return {
      title: this.title,
    };
  }

  @Watch("dialog")
  onPropertyChanged(value: string, oldValue: string) {
    value || this.close();
  }

  @Watch("dialogDelete")
  onPropertyChangedDelete(value: string, oldValue: string) {
    value || this.closeDelete();
  }

  get formTitle() {
    return this.editedIndex === -1
      ? this.$i18n.t("stores.new")
      : this.$i18n.t("stores.edit");
  }

  async mounted() {
    this.getStores();
    this.fillCountries();
    this.role = this.$auth.$state.user.Role;
  }

  addStore() {
    this.toggleSpinner(true);
    if (this.editedIndex === -1) {
      this.$axios
        .post(this.$config.internalApi + "stores", {
          countryCode: this.editedItem.country,
          url: this.editedItem.url,
          feed: this.editedItem.feed,
          window: this.editedItem.window,
          accountRefer:
            this.$auth.$state.user.Parent !== ""
              ? this.$auth.$state.user.Parent
              : this.$auth.$state.user.Id,
        })
        .then((response: any) => {
          if (response.data.success === true) {
            this.updateText(this.$i18n.t("stores.updated"));
            this.updateColor("green");
            this.updateShow(true);
            if (this.$auth.$state.user.Parent !== "") {
              this.$nuxt.$emit("refreshStores", this.$auth.$state.user.Parent);
            } else {
              this.$nuxt.$emit("refreshStores", this.$auth.$state.user.Id);
            }
            this.getStores();
            this.toggleSpinner(false);
            this.$refs.form.resetValidation();
          } else {
            this.updateText(response.data.error);
            this.updateColor("red");
            this.updateShow(true);
            this.toggleSpinner(false);
          }
        });
    } else {
      this.$axios
        .put(this.$config.internalApi + "stores", {
          countryCode: this.editedItem.country,
          url: this.editedItem.url,
          ID: this.editedItem.id,
          feed: this.editedItem.feed,
          window: this.editedItem.window
        })
        .then((response: any) => {
          if (response.data.success === true) {
            this.updateText(this.$i18n.t("stores.created"));
            this.updateColor("green");
            this.updateShow(true);
            this.getStores();
            this.toggleSpinner(false);
            this.$refs.form.resetValidation();
          } else {
            this.updateText(response.data.error);
            this.updateColor("red");
            this.updateShow(true);
            this.toggleSpinner(false);
          }
        });
    }
  }

  getStores() {
    this.dialog = false;
    const id =
      this.$auth.$state.user.Parent !== ""
        ? this.$auth.$state.user.Parent
        : this.$auth.$state.user.Id;
    this.$axios
      .get(this.$config.internalApi + "stores/" + id)
      .then((response: any) => {
        if (response.data.success === true) {
          this.stores = response.data.stores;
        } else {
          this.updateText(response.data.error);
          this.updateColor("red");
          this.updateShow(true);
        }
      });
  }

  editItem(item: StoreBO) {
    this.editedIndex = this.stores.indexOf(item);
    this.editedItem.url = item.Url;
    this.editedItem.country = item.CountryCode.toUpperCase();
    this.editedItem.id = item.Id;
    this.editedItem.feed = item.XmlFeed;
    this.editedItem.window = item.Window;
    this.dialog = true;
  }

  deleteItem(item: StoreBO) {
    this.editedIndex = this.stores.indexOf(item);
    this.editedItem.url = item.Url;
    this.editedItem.country = item.CountryCode.toUpperCase();
    this.editedItem.id = item.Id;
    this.editedItem.feed = item.XmlFeed;
    this.editedItem.window = item.Window;
    this.dialogDelete = true;
  }

  deleteItemConfirm() {
    this.toggleSpinner(true);
    this.$axios
      .delete(this.$config.internalApi + "stores/" + this.editedItem.id)
      .then((response: any) => {
        if (response.data.success) {
          this.updateText(this.$i18n.t("stores.deleted"));
          this.updateColor("green");
          this.updateShow(true);
          this.getStores();
          if (this.$auth.$state.user.Parent !== "") {
            this.$nuxt.$emit("refreshStores", this.$auth.$state.user.Parent);
          } else {
            this.$nuxt.$emit("refreshStores", this.$auth.$state.user.Id);
          }
          this.closeDelete();
          this.toggleSpinner(false);
        } else {
          this.updateText(response.data.error);
          this.updateColor("red");
          this.updateShow(true);
          this.toggleSpinner(false);
        }
      });
  }

  detailItem(item: StoreBO) {
    this.spCode = item.Code;
    this.dialogDetail = true;
  }

  close() {
    this.dialog = false;
    this.$nextTick(() => {
      this.$refs.form.resetValidation();
      this.editedItem = Object.assign(
        {},
        {
          id: "",
          url: "",
          country: "SK",
        }
      );
      this.editedIndex = -1;
    });
  }

  closeDelete() {
    this.dialogDelete = false;
    this.$nextTick(() => {
      this.editedItem = Object.assign(
        {},
        {
          id: "",
          url: "",
          country: "SK",
        }
      );
      this.editedIndex = -1;
    });
  }

  closeDetail() {
    this.dialogDetail = false;
    this.$nextTick(() => {
      this.editedItem = Object.assign(
        {},
        {
          id: "",
          url: "",
          country: "SK",
        }
      );
      this.editedIndex = -1;
    });
  }

  showSuppliers(item: StoreBO) {
    this.supplierStore = item;
    this.dialogSuppliers = true;
    this.storeName = item.Url;
    this.$axios.get("/supplier/" + item.Id).then((response: any) => {
      if (response.data.success) {
        this.suppliers = response.data.suppliers;
      } else {
        this.updateText(response.data.error);
        this.updateColor("red");
        this.updateShow(true);
      }
    });
  }

  showInvoices(item: StoreBO) {
    this.supplierStore = item;
    this.dialogInvoices = true;
    this.storeName = item.Url;
    this.$axios.get("/invoice/" + item.Id).then((response: any) => {
      if (response.data.success) {
        this.invoices = response.data.invoices;
      } else {
        this.updateText(response.data.error);
        this.updateColor("red");
        this.updateShow(true);
      }
    });
  }

  closeSuppliers() {
    this.dialogSuppliers = false;
  }

  closeInvoices() {
    this.dialogInvoices = false;
  }

  deleteSupplier(item: Supplier) {
    this.supplier = item;
    this.dialogSupplierDelete = true;
  }

  deleteSupplierConfirm() {
    this.toggleSpinner(true);
    this.$axios
      .delete("/supplier/" + this.supplier.Id)
      .then((response: any) => {
        if (response.data.success) {
          this.updateText(this.$i18n.t("suppliers.deleted"));
          this.updateColor("green");
          this.updateShow(true);
          this.showSuppliers(this.supplierStore);
          this.toggleSpinner(false);
          this.closeSupplierDelete();
        } else {
          this.updateText(response.data.error);
          this.updateColor("red");
          this.updateShow(true);
          this.toggleSpinner(false);
          this.closeSupplierDelete();
        }
      });
  }

  closeSupplierDelete() {
    this.dialogSupplierDelete = false;
  }

  closeInvoiceDelete() {
    this.dialogInvoiceDelete = false;
  }

  editSupplier(item: Supplier) {
    this.supplier = item;
    this.dialogAddSupplier = true;
  }

  editInvoice(item: Invoice) {
    this.invoice = item;
    this.dialogAddInvoice = true;
  }

  deleteInvoice(item: Invoice) {
    this.invoice = item;
    this.dialogInvoiceDelete = true;
  }

  deleteInvoiceConfirm() {
    this.toggleSpinner(true);
    this.$axios
      .delete("/invoice/" + this.supplier.Id)
      .then((response: any) => {
        if (response.data.success) {
          this.updateText(this.$i18n.t("invoices.deleted"));
          this.updateColor("green");
          this.updateShow(true);
          this.showInvoices(this.supplierStore);
          this.toggleSpinner(false);
          this.closeSupplierDelete();
        } else {
          this.updateText(response.data.error);
          this.updateColor("red");
          this.updateShow(true);
          this.toggleSpinner(false);
          this.closeSupplierDelete();
        }
      });
  }

  get formTitleSupplier() {
    return this.supplier.Id === ""
      ? this.$i18n.t("stores.new")
      : this.$i18n.t("stores.edit");
  }

  closeAddSuppliers() {
    this.dialogAddSupplier = false;
    this.supplier = {
      Id: "",
      Name: "",
      Street: "",
      City: "",
      Zip: "",
      Email: "",
      Phone: "",
      Person: "",
      Country: "",
      Template: "",
      Subject: "",
    };
    this.template = "";
    this.$refs.supplierForm.resetValidation();
  }

  addSupplier() {
    this.toggleSpinner(true);
    if (this.supplier.Id !== "") {
      this.$axios
        .put("/supplier", {
          Id: this.supplier.Id,
          Name: this.supplier.Name,
          Street: this.supplier.Street,
          City: this.supplier.City,
          Zip: this.supplier.Zip,
          Email: this.supplier.Email,
          Phone: this.supplier.Phone,
          Person: this.supplier.Person,
          Template: this.supplier.Template,
          Subject: this.supplier.Subject,
        })
        .then((response: any) => {
          if (response.data.success) {
            this.updateText(this.$i18n.t("suppliers.updated"));
            this.updateColor("green");
            this.updateShow(true);
            this.closeAddSuppliers();
            this.toggleSpinner(false);
          } else {
            this.updateText(response.data.error);
            this.updateColor("red");
            this.updateShow(true);
            this.toggleSpinner(false);
          }
        });
    } else {
      this.$axios
        .post("/supplier", {
          Name: this.supplier.Name,
          Street: this.supplier.Street,
          City: this.supplier.City,
          Zip: this.supplier.Zip,
          Email: this.supplier.Email,
          Phone: this.supplier.Phone,
          Person: this.supplier.Person,
          StoreRefer: this.supplierStore.Id,
          Template: this.supplier.Template,
          Subject: this.supplier.Subject,
        })
        .then((response: any) => {
          if (response.data.success) {
            this.suppliers.push(this.supplier);
            this.closeAddSuppliers();
            this.updateText(this.$i18n.t("suppliers.created"));
            this.updateColor("green");
            this.updateShow(true);
            this.toggleSpinner(false);
            this.showSuppliers(this.supplierStore);
          } else {
            this.updateText(response.data.error);
            this.updateColor("red");
            this.updateShow(true);
            this.toggleSpinner(false);
          }
        });
    }
  }

  onEditorBlur(quill: any) {}

  onEditorFocus(quill: any) {}

  onEditorReady(quill: any) {}

  onEditorChange(quill: any, html: any, text: any) {}

  get formTitleInvoice() {
    return this.invoice.Id === ""
      ? this.$i18n.t("invoices.new")
      : this.$i18n.t("invoices.edit");
  }

  closeAddInvoice() {
    this.dialogAddInvoice = false;
    this.invoice = {
      Id: "",
      DueDate: "",
      Amount: 0,
      Currency: "",
    };
    this.$refs.invoiceForm.resetValidation();
  }

  addInvoice() {
    this.toggleSpinner(true);
    if (this.invoice.Id !== "") {
      this.$axios
        .put("/invoice", {
          Id: this.invoice.Id,
          DueDate: new Date(this.invoice.DueDate).toISOString(),
          Amount: this.invoice.Amount,
          Currency: this.invoice.Currency,
        })
        .then((response: any) => {
          if (response.data.success) {
            this.closeAddInvoice();
            this.updateText(this.$i18n.t("invoices.updated"));
            this.updateColor("green");
            this.updateShow(true);
            this.toggleSpinner(false);
          } else {
            this.updateText(response.data.error);
            this.updateColor("red");
            this.updateShow(true);
            this.toggleSpinner(false);
          }
        });
    } else {
      this.$axios
        .post("/invoice", {
          DueDate: new Date(this.invoice.DueDate).toISOString(),
          Amount: this.invoice.Amount,
          Currency: this.invoice.Currency,
          StoreRefer: this.supplierStore.Id,
        })
        .then((response: any) => {
          if (response.data.success) {
            this.invoices.push(this.invoice);
            this.closeAddInvoice();
            this.updateText(this.$i18n.t("invoices.created"));
            this.updateColor("green");
            this.updateShow(true);
            this.toggleSpinner(false);
            this.showInvoices(this.supplierStore);
          } else {
            this.updateText(response.data.error);
            this.updateColor("red");
            this.updateShow(true);
            this.toggleSpinner(false);
          }
        });
    }
  }

  customFilter(item: ICountry, queryText: string, itemText: string) {
    const textOne = item.Name.toLowerCase();
    const textTwo = item.Code.toLowerCase();
    const searchText = queryText.toLowerCase();

    return textOne.indexOf(searchText) > -1 || textTwo.indexOf(searchText) > -1;
  }

  fillCountries() {
    this.$axios
      .get(
        window.location.protocol +
          "//" +
          window.location.host +
          "/countries.json"
      )
      .then((res: any) => (this.countries = res.data));
  }
}
</script>

<style scoped>
.v-icon {
  color: #000059;
}

.theme--dark.v-icon {
  color: #e1863b !important;
}

.theme--dark.v-application .primary {
  border-color: #e1863b !important;
  background-color: #e1863b !important;
}
.theme--light.v-card > .v-card__text,
.theme--light.v-card > .v-card__subtitle {
  color: #fff !important;
}
.stores .v-label.theme--light {
  color: #fff !important;
}
.theme--light .v-card__title,
.v-dialog > .v-card > .v-card__title {
  background: #4633af !important;
  color: #fff !important;
}
.v-dialog .container {
  color: #000 !important;
}
.gray {
  background: black;
}

.theme--dark .v-dialog .v-card__title {
  background-color: #4f5249 !important;
}
</style>
