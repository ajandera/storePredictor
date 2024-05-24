<template>
  <v-container fluid>
    <v-row justify="center" align="center">
      <v-col cols="12">
        <v-data-table
          :headers="headers"
          :items="accounts"
          :items-per-page="10"
          class="elevation-1"
        >
          <template v-slot:top>
            <v-toolbar flat>
              <v-toolbar-title>{{ $t("users.title") }}</v-toolbar-title>
              <v-divider class="mx-4" inset vertical></v-divider>
              <v-spacer></v-spacer>
              <v-dialog v-model="dialog" max-width="500px">
                <template v-slot:activator="{ on, attrs }">
                  <Help
                    title="User Management"
                    keyString="help-userManagement"
                  />
                  <v-btn
                    color="primary"
                    dark
                    class="mb-2 mr-2"
                    v-bind="attrs"
                    v-on="on"
                  >
                    {{ $t("users.add") }}
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
                          <v-col cols="12">
                            <v-text-field
                              v-model="editedItem.name"
                              :counter="50"
                              :rules="nameRules"
                              :label="$t('users.name')"
                              required
                            ></v-text-field>
                          </v-col>
                          <v-col cols="12">
                            <v-text-field
                              v-model="editedItem.email"
                              :counter="50"
                              :rules="emailRules"
                              :label="$t('users.email')"
                              required
                            ></v-text-field>
                          </v-col>
                          <v-col cols="12">
                            <v-select
                              v-model="editedItem.role"
                              :items="roles"
                              :rules="[(v) => !!v || 'Item is required']"
                              :label="$t('users.roles')"
                              required
                            ></v-select>
                          </v-col>
                          <v-col cols="12">
                            <v-autocomplete
                              :items="countries"
                              :filter="customFilter"
                              item-text="Name"
                              item-value="Code"
                              :rules="[(v) => !!v || 'Item is required']"
                              :label="$t('stores.form.country')"
                              v-model="editedItem.countryCode"
                            ></v-autocomplete>
                          </v-col>
                          <v-col cols="12">
                            <v-text-field
                              v-model="password"
                              type="password"
                              :rules="passwordRules"
                              :label="$t('profile.password')"
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
                        @click="addUser"
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
                    $t("users.delete")
                  }}</v-card-title>
                  <v-card-actions>
                    <v-spacer></v-spacer>
                    <v-btn color="blue darken-1" text @click="closeDelete">{{
                      $t("stores.cancel")
                    }}</v-btn>
                    <v-btn
                      color="blue darken-1"
                      text
                      @click="deleteItemConfirm"
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
                  <v-btn v-bind="attrs" v-on="on" @click="editItem(item)">
                    <v-icon small class="mr-2" color="secondary">
                      mdi-pencil
                    </v-icon>
                  </v-btn>
                </template>
                <span>{{ $t("tooltips.edit") }}</span>
              </v-tooltip>
              <v-tooltip bottom>
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
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import { Component, Vue, Watch } from "nuxt-property-decorator";
import Account from "~/model/Account";
import User from "~/model/User";
import { namespace } from "vuex-class";
import ICountry from "~/model/ICountry";
const snackbar = namespace("Snackbar");
const spinner = namespace("Spinner");

@Component
export default class UsersPage extends Vue {
  @snackbar.Action
  public updateText!: (newText: string) => void;

  @snackbar.Action
  public updateColor!: (newColor: string) => void;

  @snackbar.Action
  public updateShow!: (newShow: boolean) => void;

  @spinner.Action
  public toggleSpinner!: (newShow: boolean) => void;

  $i18n: any;
  accounts: Account[] = [];
  roles: Array<string> = ["Admin", "Viewer"];
  dialog: boolean = false;
  dialogDelete: boolean = false;
  password: string = "";
  $refs: any;
  headers: Array<any> = [
    {
      // @ts-ignore
      text: this.$i18n.t("users.name"),
      align: "start",
      sortable: true,
      value: "Name",
    },
    // @ts-ignore
    { text: this.$i18n.t("users.email"), value: "Email" },
    // @ts-ignore
    { text: this.$i18n.t("users.roles"), value: "Role", sortable: true },
    // @ts-ignore
    { text: this.$i18n.t("stores.actions"), value: "actions", sortable: false },
  ];
  editedIndex: number = -1;
  editedItem: User = {
    id: "",
    name: "",
    role: "Viewer",
    email: "",
    countryCode: "SK",
  };
  // @ts-ignore
  title: string = this.$i18n.t("users.title");
  valid: boolean = true;
  nameRules: Array<any> = [
    (v: any) => !!v || this.$i18n.t("validation.name"),
    (v: string | any[]) =>
      (v && v.length <= 50) || this.$i18n.t("validation.name"),
  ];
  emailRules: Array<any> = [
    (v: any) => !!v || this.$i18n.t("profile.emailis"),
    (v: string) => /.+@.+\..+/.test(v) || this.$i18n.t("profile.emailmust"),
  ];
  passwordRules: Array<any> = [];
  countries: Array<ICountry> = [];
  $axios: any;
  countryCode: any;
  $t: any;

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

  mounted() {
    this.getUsers();
    this.fillCountries();
  }

  addUser() {
    this.toggleSpinner(true);
    if (this.editedIndex === -1) {
      this.$axios
        .get(
          this.$config.internalApi + "account/email/" + this.editedItem.email
        )
        .then((response: any) => {
          if (
            response.data.success === true &&
            response.data.account.Id === "00000000-0000-0000-0000-000000000000"
          ) {
            this.$axios
              .put(this.$config.internalApi + "accounts", {
                email: this.editedItem.email,
                name: this.editedItem.name,
                countryCode: this.editedItem.countryCode,
                password: this.password,
                role: this.editedItem.role,
                parent: this.$auth.$state.user.Id,
                planRefer: this.$auth.$state.user.PlanRefer
              })
              .then((response: any) => {
                if (response.data.success === true) {
                  this.getUsers();
                  this.updateText(this.$i18n.t("users.created"));
                  this.updateColor("green");
                  this.updateShow(true);
                  this.toggleSpinner(false);
                } else {
                  this.updateText(response.data.error);
                  this.updateColor("red");
                  this.updateShow(true);
                  this.toggleSpinner(false);
                }
                this.$refs.form.resetValidation();
              });
          } else {
            this.updateText(this.$i18n.t("users.already"));
            this.updateColor("red");
            this.updateShow(true);
            this.toggleSpinner(false);
          }
        });
    } else {
      this.$axios
        .put(this.$config.internalApi + "accounts", {
          Id: this.editedItem.id,
          name: this.editedItem.name,
          email: this.editedItem.email,
          countryCode: this.countryCode,
          password: this.password,
          role: this.editedItem.role,
        })
        .then((response: any) => {
          if (response.data.success === true) {
            this.getUsers();
            this.updateText(this.$i18n.t("users.updated"));
            this.updateColor("green");
            this.updateShow(true);
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

  getUsers() {
    this.dialog = false;
    this.$axios
      .get("/account/child/" + this.$auth.$state.user.Id)
      .then((response: any) => {
        if (response.data.success === true) {
          this.accounts = response.data.accounts;
        } else {
          this.updateText(response.data.error);
          this.updateColor("red");
          this.updateShow(true);
        }
      });
  }

  editItem(item: Account) {
    this.editedIndex = this.accounts.indexOf(item);
    this.editedItem.name = item.Name;
    this.editedItem.email = item.Email;
    this.editedItem.role = item.Role;
    this.editedItem.countryCode = item.CountryCode;
    this.editedItem.id = item.Id;
    this.dialog = true;
  }

  deleteItem(item: Account) {
    this.editedIndex = this.accounts.indexOf(item);
    this.editedItem.name = item.Name;
    this.editedItem.email = item.Email;
    this.editedItem.role = item.Role;
    this.editedItem.id = item.Id;
    this.dialogDelete = true;
  }

  deleteItemConfirm() {
    this.toggleSpinner(true);
    this.$axios
      .delete(this.$config.internalApi + "accounts/" + this.editedItem.id)
      .then((response: any) => {
        if (response.data.success === true) {
          this.getUsers();
          this.closeDelete();
          this.toggleSpinner(false);
        } else {
          this.updateText(response.data.error);
          this.updateColor("red");
          this.updateShow(true);
        }
      });
  }

  close() {
    this.dialog = false;
    this.$nextTick(() => {
      this.$refs.form.resetValidation();
      this.editedItem = Object.assign(
        {},
        { id: "", name: "", role: "Viewer", email: "", countryCode: "SK" }
      );
      this.editedIndex = -1;
    });
  }

  closeDelete() {
    this.dialogDelete = false;
    this.$nextTick(() => {
      this.editedItem = Object.assign(
        {},
        { id: "", name: "", role: "Viewer", email: "", countryCode: "SK" }
      );
      this.editedIndex = -1;
    });
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
.v-dialog .v-card__title {
  background: #4633af !important;
  color: #fff !important;
}
.theme--dark .v-dialog .v-card__title {
  background-color: #4f5249 !important;
}
.v-dialog .container {
  color: #000 !important;
}
</style>
