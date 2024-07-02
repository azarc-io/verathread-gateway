import { CodegenConfig } from "@graphql-codegen/cli";

const config: CodegenConfig = {
  schema: 'http://localhost:6010/graphql',
  documents: ['src/**/*.graphqls', 'src/*.graphqls', 'src/**/*.vue'],
  generates: {
    "./src/gql/": {
      preset: "client",
      config: {
        useTypeImports: true
      }
    },
  },
};

export default config;
