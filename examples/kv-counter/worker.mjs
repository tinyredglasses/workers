import "../assets/polyfill_performance.js";
import "../assets/wasm_exec.js";
import mod from "./dist/app.wasm";

const go = new Go();

const load = WebAssembly.instantiate(mod, go.importObject).then((instance) => {
  go.run(instance);
  return instance;
});

async function processRequest(event) {
  const req = event.request;
  await load;
  return handleRequest(req);
}

addEventListener("fetch", (event) => {
  event.respondWith(processRequest(event));
})
