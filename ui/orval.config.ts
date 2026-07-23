import { defineConfig } from "orval"

export default defineConfig({
  depositum: {
    input: {
      target: "../docs/swagger.json",
    },

    output: {
      mode: "tags-split",
      target: "./app/lib/api/generated",
      schemas: "./app/lib/api/model",

      client: "swr",
      httpClient: "fetch",
      baseUrl: "/api/v1",
    },
  },
})
