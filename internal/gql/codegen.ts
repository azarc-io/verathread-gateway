import type {CodegenConfig} from '@graphql-codegen/cli';

const config: CodegenConfig = {
    schema: ["schema/**/*.graphqls"],
    generates: {
        'schema.graphql': {
            plugins: ['schema-ast'],
            config: {
                includeDirectives: true
            },
        },
    },
};
export default config;
