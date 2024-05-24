// store/modules/snackbar.ts
import { VuexModule, Module, Mutation, Action } from "vuex-module-decorators";

@Module({ namespaced: true })
class Snackbar extends VuexModule {
  public text: string = "";
  public color: string = "";
  public show: boolean = false;

  @Mutation
  public setText(newText: string): void {
    this.text = newText;
  }

  @Mutation
  public setColor(newColor: string): void {
    this.color = newColor;
  }

  @Mutation
  public setShow(newVisibility: boolean): void {
    this.show = newVisibility;
  }

  @Action({ rawError: true })
  public updateText(newText: string): void {
    this.context.commit("setText", newText);
  }

  @Action({ rawError: true })
  public updateColor(newColor: string): void {
    this.context.commit("setColor", newColor);
  }

  @Action({ rawError: true })
  public updateShow(newVisibility: boolean): void {
    this.context.commit("setShow", newVisibility);
  }
}
export default Snackbar;
