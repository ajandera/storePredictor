import { defineComponent, h, ref, watch, onMounted } from "vue";
import VueApexCharts from "vue3-apexcharts";
import mermaid from "mermaid";

export default defineNuxtPlugin((nuxtApp) => {
  // 1. <apexchart>
  nuxtApp.vueApp.component("apexchart", VueApexCharts);

  // 2. <country-flag :country="code" size="small" />
  const CountryFlag = defineComponent({
    name: "CountryFlag",
    props: {
      country: { type: String, default: "us" },
      size: { type: String, default: "normal" },
    },
    setup(props) {
      const toFlagEmoji = (code: string) => {
        const clean = (code || "us")
          .trim()
          .toLowerCase()
          .replace(/^en$/, "us")
          .replace(/^cs$/, "cz")
          .slice(0, 2)
          .toUpperCase();
        if (clean.length !== 2) return "🏳️";
        return String.fromCodePoint(
          ...Array.from(clean).map((c) => 127397 + c.charCodeAt(0))
        );
      };
      return () =>
        h(
          "span",
          {
            class: ["country-flag", `country-flag--${props.size}`],
            style: { marginRight: "6px", display: "inline-block" },
          },
          toFlagEmoji(props.country)
        );
    },
  });
  nuxtApp.vueApp.component("country-flag", CountryFlag);

  // 3. <quill-editor v-model="content" :options="editorOption" @blur @focus @ready />
  const QuillEditor = defineComponent({
    name: "QuillEditor",
    props: {
      modelValue: { type: String, default: "" },
      value: { type: String, default: "" },
      options: { type: Object, default: () => ({}) },
    },
    emits: ["update:modelValue", "input", "blur", "focus", "ready", "change"],
    setup(props, { emit }) {
      onMounted(() => {
        emit("ready", {});
      });
      return () =>
        h("textarea", {
          class: "quill-editor-textarea",
          value: props.modelValue || props.value || "",
          rows: 6,
          style: {
            width: "100%",
            padding: "10px",
            border: "1px solid #ccc",
            borderRadius: "4px",
            backgroundColor: "#fff",
            color: "#000",
          },
          onInput: (e: Event) => {
            const val = (e.target as HTMLTextAreaElement).value;
            emit("update:modelValue", val);
            emit("input", val);
            emit("change", { html: val, text: val });
          },
          onBlur: (e: Event) => emit("blur", e),
          onFocus: (e: Event) => emit("focus", e),
        });
    },
  });
  nuxtApp.vueApp.component("quill-editor", QuillEditor);

  // 4. <vue-mermaid-string :value="diagram" />
  try {
    mermaid.initialize({ startOnLoad: false, securityLevel: "loose" });
  } catch {
    // ignore init error
  }

  const VueMermaidString = defineComponent({
    name: "VueMermaidString",
    props: {
      value: { type: String, default: "" },
    },
    setup(props) {
      const svgHtml = ref("");
      const renderDiagram = async () => {
        if (!props.value) {
          svgHtml.value = "";
          return;
        }
        try {
          const id = "mermaid-" + Math.random().toString(36).slice(2, 10);
          const { svg } = await mermaid.render(id, props.value);
          svgHtml.value = svg;
        } catch {
          svgHtml.value = `<pre>${props.value}</pre>`;
        }
      };
      onMounted(renderDiagram);
      watch(() => props.value, renderDiagram);
      return () =>
        h("div", {
          class: "vue-mermaid-string",
          innerHTML: svgHtml.value,
        });
    },
  });
  nuxtApp.vueApp.component("vue-mermaid-string", VueMermaidString);

  // 5. <paypal-checkout />
  const PaypalCheckout = defineComponent({
    name: "PaypalCheckout",
    props: {
      amount: { type: [String, Number], default: "0" },
      currency: { type: String, default: "EUR" },
      env: { type: String, default: "sandbox" },
      client: { type: Object, default: () => ({}) },
      invoiceNumber: { type: String, default: "" },
      braintree: { type: null, default: null },
    },
    setup(props) {
      return () =>
        h(
          "div",
          { class: "paypal-checkout-wrapper mt-5" },
          `PayPal (${props.amount} ${props.currency})`
        );
    },
  });
  nuxtApp.vueApp.component("paypal-checkout", PaypalCheckout);
});
