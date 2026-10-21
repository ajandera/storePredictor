import { storeState, storeActions, pendingMeta } from "./store";

function extractPropName(keyOrContext: any): string {
  if (typeof keyOrContext === "string") return keyOrContext;
  if (keyOrContext && typeof keyOrContext === "object" && "name" in keyOrContext) {
    return String(keyOrContext.name);
  }
  return "";
}

export function namespace(ns: "Snackbar" | "Spinner") {
  return {
    Action(targetOrValue: any, keyOrContext?: any) {
      const propertyKey = extractPropName(keyOrContext);
      if (!propertyKey) return;
      const fn = function (...args: any[]) {
        return storeActions[ns]?.[propertyKey]?.(...args);
      };
      pendingMeta.vuexActions[propertyKey] = fn;
      if (targetOrValue && typeof targetOrValue === "object" && targetOrValue.constructor) {
        const ctor = targetOrValue.constructor;
        if (!ctor.__vuexActions) ctor.__vuexActions = {};
        ctor.__vuexActions[propertyKey] = fn;
      }
    },
    State(targetOrValue: any, keyOrContext?: any) {
      const propertyKey = extractPropName(keyOrContext);
      if (!propertyKey) return;
      const getter = function () {
        return (storeState as any)[ns]?.[propertyKey];
      };
      pendingMeta.vuexStates[propertyKey] = getter;
      if (targetOrValue && typeof targetOrValue === "object" && targetOrValue.constructor) {
        const ctor = targetOrValue.constructor;
        if (!ctor.__vuexStates) ctor.__vuexStates = {};
        ctor.__vuexStates[propertyKey] = getter;
      }
    },
  };
}
