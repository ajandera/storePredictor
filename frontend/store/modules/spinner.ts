// store/modules/spinner.ts
import { VuexModule, Module, Mutation, Action } from "vuex-module-decorators";

@Module({ namespaced: true })
class Spinner extends VuexModule {
  public show: boolean = false;

  @Mutation
  public setShow(newVisibility: boolean): void {
    this.show = newVisibility;
  }

  @Action({ rawError: true })
  public toggleSpinner(newVisibility: boolean): void {
    this.context.commit("setShow", newVisibility);
  }
}
export default Spinner;
