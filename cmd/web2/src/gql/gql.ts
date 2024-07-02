/* eslint-disable */
import * as types from './graphql';
import type { TypedDocumentNode as DocumentNode } from '@graphql-typed-document-node/core';

/**
 * Map of all GraphQL operations in the project.
 *
 * This map has several performance disadvantages:
 * 1. It is not tree-shakeable, so it will include all operations in the project.
 * 2. It is not minifiable, so the string of a GraphQL query will be multiple times inside the bundle.
 * 3. It does not support dead code elimination, so it will add unused operations.
 *
 * Therefore it is highly recommended to use the babel or swc plugin for production.
 */
const documents = {
    "query FetchShellConfig($tenant: String!) {\n  shellConfiguration(tenantId: $tenant) {\n    defaultRoute\n    categories {\n      category\n      priority\n      title\n      entries {\n        id\n        title\n        subTitle\n        authRequired\n        healthy\n        hidden\n        icon\n        module {\n          exposedModule\n          moduleName\n          outlet\n          path\n          remoteEntry\n        }\n      }\n    }\n    slots {\n      authRequired\n      priority\n      description\n      slot\n      module {\n        remoteEntry\n        path\n        moduleName\n        exposedModule\n      }\n    }\n  }\n}\n\nsubscription SubscribeToShellConfig($tenant: String!, $events: [ShellConfigEventType!]!) {\n  shellConfiguration(tenantId: $tenant, events: $events) {\n    configuration {\n      defaultRoute\n      categories {\n        category\n        priority\n        title\n        entries {\n          id\n          title\n          subTitle\n          authRequired\n          healthy\n          hidden\n          icon\n          module {\n            exposedModule\n            moduleName\n            outlet\n            path\n            remoteEntry\n          }\n        }\n      }\n      slots {\n        authRequired\n        priority\n        description\n        slot\n        module {\n          remoteEntry\n          path\n          moduleName\n          exposedModule\n        }\n      }\n    }\n    eventType\n  }\n}": types.FetchShellConfigDocument,
};

/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 *
 *
 * @example
 * ```ts
 * const query = graphql(`query GetUser($id: ID!) { user(id: $id) { name } }`);
 * ```
 *
 * The query argument is unknown!
 * Please regenerate the types.
 */
export function graphql(source: string): unknown;

/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "query FetchShellConfig($tenant: String!) {\n  shellConfiguration(tenantId: $tenant) {\n    defaultRoute\n    categories {\n      category\n      priority\n      title\n      entries {\n        id\n        title\n        subTitle\n        authRequired\n        healthy\n        hidden\n        icon\n        module {\n          exposedModule\n          moduleName\n          outlet\n          path\n          remoteEntry\n        }\n      }\n    }\n    slots {\n      authRequired\n      priority\n      description\n      slot\n      module {\n        remoteEntry\n        path\n        moduleName\n        exposedModule\n      }\n    }\n  }\n}\n\nsubscription SubscribeToShellConfig($tenant: String!, $events: [ShellConfigEventType!]!) {\n  shellConfiguration(tenantId: $tenant, events: $events) {\n    configuration {\n      defaultRoute\n      categories {\n        category\n        priority\n        title\n        entries {\n          id\n          title\n          subTitle\n          authRequired\n          healthy\n          hidden\n          icon\n          module {\n            exposedModule\n            moduleName\n            outlet\n            path\n            remoteEntry\n          }\n        }\n      }\n      slots {\n        authRequired\n        priority\n        description\n        slot\n        module {\n          remoteEntry\n          path\n          moduleName\n          exposedModule\n        }\n      }\n    }\n    eventType\n  }\n}"): (typeof documents)["query FetchShellConfig($tenant: String!) {\n  shellConfiguration(tenantId: $tenant) {\n    defaultRoute\n    categories {\n      category\n      priority\n      title\n      entries {\n        id\n        title\n        subTitle\n        authRequired\n        healthy\n        hidden\n        icon\n        module {\n          exposedModule\n          moduleName\n          outlet\n          path\n          remoteEntry\n        }\n      }\n    }\n    slots {\n      authRequired\n      priority\n      description\n      slot\n      module {\n        remoteEntry\n        path\n        moduleName\n        exposedModule\n      }\n    }\n  }\n}\n\nsubscription SubscribeToShellConfig($tenant: String!, $events: [ShellConfigEventType!]!) {\n  shellConfiguration(tenantId: $tenant, events: $events) {\n    configuration {\n      defaultRoute\n      categories {\n        category\n        priority\n        title\n        entries {\n          id\n          title\n          subTitle\n          authRequired\n          healthy\n          hidden\n          icon\n          module {\n            exposedModule\n            moduleName\n            outlet\n            path\n            remoteEntry\n          }\n        }\n      }\n      slots {\n        authRequired\n        priority\n        description\n        slot\n        module {\n          remoteEntry\n          path\n          moduleName\n          exposedModule\n        }\n      }\n    }\n    eventType\n  }\n}"];

export function graphql(source: string) {
  return (documents as any)[source] ?? {};
}

export type DocumentType<TDocumentNode extends DocumentNode<any, any>> = TDocumentNode extends DocumentNode<  infer TType,  any>  ? TType  : never;