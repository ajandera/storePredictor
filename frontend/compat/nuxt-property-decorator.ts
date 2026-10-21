import { defineComponent } from "vue";
import { useHead } from "#imports";
import { pendingMeta } from "./store";

let activeVm: any = null;

export class Vue {
  $i18n: any;
  $t: any;
  $axios: any;
  $auth: any;
  $config: any;
  $route: any;
  $router: any;
  $nuxt: any;
  $refs: any;
  $store: any;
  $vuetify: any;
  $nextTick: any;

  constructor() {
    if (activeVm) {
      const vm = activeVm;
      return new Proxy(this, {
        get(target: any, prop: string | symbol, receiver: any) {
          if (Reflect.has(target, prop)) {
            return Reflect.get(target, prop, receiver);
          }
          if (vm && prop in vm) {
            const v = vm[prop];
            return typeof v === "function" ? v.bind(vm) : v;
          }
          return undefined;
        },
        set(target: any, prop: string | symbol, value: any) {
          target[prop] = value;
          return true;
        },
      });
    }
  }
}

function extractPropName(keyOrContext: any): string {
  if (typeof keyOrContext === "string") return keyOrContext;
  if (keyOrContext && typeof keyOrContext === "object" && "name" in keyOrContext) {
    return String(keyOrContext.name);
  }
  return "";
}

export function Prop(options: any = {}) {
  return function (targetOrValue: any, keyOrContext?: any) {
    const propertyKey = extractPropName(keyOrContext);
    if (!propertyKey) return;
    pendingMeta.props[propertyKey] = options || {};
    if (targetOrValue && typeof targetOrValue === "object" && targetOrValue.constructor) {
      const ctor = targetOrValue.constructor;
      if (!ctor.__props) ctor.__props = {};
      ctor.__props[propertyKey] = options || {};
    }
  };
}

export function Watch(path: string, options: any = {}) {
  return function (targetOrValue: any, keyOrContext?: any) {
    const propertyKey = extractPropName(keyOrContext);
    if (!propertyKey) return;
    pendingMeta.watchers.push({ path, handler: propertyKey, options });
    if (targetOrValue && typeof targetOrValue === "object" && targetOrValue.constructor) {
      const ctor = targetOrValue.constructor;
      if (!ctor.__watchers) ctor.__watchers = [];
      ctor.__watchers.push({ path, handler: propertyKey, options });
    }
  };
}

const LIFECYCLE_HOOKS = new Set([
  "beforeCreate",
  "created",
  "beforeMount",
  "mounted",
  "beforeUpdate",
  "updated",
  "beforeUnmount",
  "unmounted",
  "destroyed",
  "beforeDestroy",
]);

function normalizeDataTableHeaders(val: any): any {
  if (
    Array.isArray(val) &&
    val.length > 0 &&
    val[0] &&
    typeof val[0] === "object" &&
    ("text" in val[0] || "value" in val[0])
  ) {
    return val.map((h: any) => ({
      ...h,
      title: h.title ?? h.text ?? "",
      key: h.key ?? h.value ?? "",
    }));
  }
  return val;
}

function buildComponent(Target: any, options: Record<string, any> = {}) {
  const proto = Target.prototype;
  const methods: Record<string, Function> = {
    ...(options.methods || {}),
    ...pendingMeta.vuexActions,
    ...(Target.__vuexActions || {}),
  };
  const computed: Record<string, any> = {
    ...(options.computed || {}),
    ...pendingMeta.vuexStates,
    ...(Target.__vuexStates || {}),
  };
  const props: Record<string, any> = {
    ...(options.props || {}),
    ...pendingMeta.props,
    ...(Target.__props || {}),
  };
  const watch: Record<string, any> = { ...(options.watch || {}) };
  const hooks: Record<string, Function> = {};

  const allWatchers = [
    ...pendingMeta.watchers,
    ...(Array.isArray(Target.__watchers) ? Target.__watchers : []),
  ];
  pendingMeta.reset();

  const descriptors = Object.getOwnPropertyDescriptors(proto);
  for (const [key, desc] of Object.entries(descriptors)) {
    if (key === "constructor") continue;
    if (desc.get || desc.set) {
      computed[key] = {
        get: desc.get,
        set: desc.set,
      };
    } else if (typeof desc.value === "function") {
      if (LIFECYCLE_HOOKS.has(key)) {
        const mappedKey =
          key === "destroyed"
            ? "unmounted"
            : key === "beforeDestroy"
            ? "beforeUnmount"
            : key;
        hooks[mappedKey] = desc.value;
      } else {
        methods[key] = desc.value;
      }
    }
  }

  for (const w of allWatchers) {
    watch[w.path] = {
      handler(this: any, newVal: any, oldVal: any) {
        return this[w.handler]?.(newVal, oldVal);
      },
      ...w.options,
    };
  }

  const userCreated = hooks.created;
  hooks.created = function (this: any) {
    if (typeof this.head === "function") {
      try {
        const headObj = this.head();
        if (headObj) useHead(headObj);
      } catch {
        // ignore head errors outside setup
      }
    }
    if (userCreated) {
      return userCreated.call(this);
    }
  };

  return defineComponent({
    name: Target.name || options.name || "ClassComponent",
    ...options,
    props,
    computed,
    methods,
    watch,
    ...hooks,
    data(vm: any) {
      const prevVm = activeVm;
      activeVm = vm || this;
      let instance: any;
      try {
        instance = new Target();
      } finally {
        activeVm = prevVm;
      }
      const dataObj: Record<string, any> = {};
      for (const key of Object.keys(instance)) {
        if (
          key.startsWith("$") ||
          key in props ||
          key in computed ||
          key in methods
        ) {
          continue;
        }
        const val = instance[key];
        if (val !== undefined) {
          dataObj[key] = normalizeDataTableHeaders(val);
        } else {
          dataObj[key] = undefined;
        }
      }
      return dataObj;
    },
  });
}

export function Component(optionsOrTarget?: any, maybeContext?: any): any {
  if (typeof optionsOrTarget === "function") {
    return buildComponent(optionsOrTarget, {});
  }
  return function (Target: any) {
    return buildComponent(Target, optionsOrTarget || {});
  };
}
