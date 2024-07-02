import { CodegenConfig } from "@graphql-codegen/cli";

const config: CodegenConfig = {
  schema: 'http://localhost:6010/graphql',
  documents: ["src/**/*.graphqls", "src/*.graphqls"],
  generates: {
    "./src/__generated__/": {
      preset: "client",
    },
    "schema.graphql": {
      plugins: ['schema-ast'],
      config: {
        includeDirectives: true,
        includeIntrospectionTypes: true
      }
    }
  },
};

export default config;
