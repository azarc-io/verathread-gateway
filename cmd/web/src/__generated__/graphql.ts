/* eslint-disable */
import { TypedDocumentNode as DocumentNode } from '@graphql-typed-document-node/core';
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
export type MakeEmpty<T extends { [key: string]: unknown }, K extends keyof T> = { [_ in K]?: never };
export type Incremental<T> = T | { [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never };
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string; output: string; }
  String: { input: string; output: string; }
  Boolean: { input: boolean; output: boolean; }
  Int: { input: number; output: number; }
  Float: { input: number; output: number; }
  Any: { input: any; output: any; }
  Duration: { input: any; output: any; }
  Time: { input: any; output: any; }
};

export type KeepAliveAppInput = {
  pkg: Scalars['String']['input'];
  version: Scalars['String']['input'];
};

export type KeepAliveAppOutput = {
  __typename?: 'KeepAliveAppOutput';
  ok: Scalars['Boolean']['output'];
  registrationRequired: Scalars['Boolean']['output'];
};

export type Mutation = {
  __typename?: 'Mutation';
  keepAlive: KeepAliveAppOutput;
  registerApp: RegisterAppOutput;
};


export type MutationKeepAliveArgs = {
  input?: InputMaybe<KeepAliveAppInput>;
};


export type MutationRegisterAppArgs = {
  input: RegisterAppInput;
};

export type Page = {
  limit?: InputMaybe<Scalars['Int']['input']>;
  page?: InputMaybe<Scalars['Int']['input']>;
};

export type PageInfo = {
  __typename?: 'PageInfo';
  next: Scalars['Int']['output'];
  page: Scalars['Int']['output'];
  perPage: Scalars['Int']['output'];
  prev: Scalars['Int']['output'];
  total: Scalars['Int']['output'];
  totalPage: Scalars['Int']['output'];
};

export type Query = {
  __typename?: 'Query';
  registeredApps?: Maybe<RegisteredAppsPage>;
  shellConfiguration: ShellConfiguration;
};


export type QueryRegisteredAppsArgs = {
  page: Page;
  sort?: InputMaybe<RegisteredAppsSort>;
  where?: InputMaybe<RegisteredAppsWhereRules>;
};


export type QueryShellConfigurationArgs = {
  tenantId: Scalars['String']['input'];
};

export enum QueryCondition {
  And = 'And',
  Or = 'Or'
}

export type QueryOperatorAndDate = {
  op: QueryOperators;
  value?: InputMaybe<Scalars['Time']['input']>;
};

export type QueryOperatorAndValue = {
  op: QueryOperators;
  value?: InputMaybe<Scalars['Any']['input']>;
};

export type QueryOperatorFieldAndValue = {
  exists?: InputMaybe<Scalars['Boolean']['input']>;
  field: Scalars['String']['input'];
  op: QueryOperators;
  value?: InputMaybe<Scalars['Any']['input']>;
};

export enum QueryOperators {
  Contains = 'Contains',
  Equal = 'Equal',
  GreaterThan = 'GreaterThan',
  GreaterThanOrEqual = 'GreaterThanOrEqual',
  In = 'In',
  LessThan = 'LessThan',
  LessThanOrEqual = 'LessThanOrEqual',
  NotEqual = 'NotEqual',
  NotIn = 'NotIn',
  Regex = 'Regex'
}

export enum QueryType {
  Date = 'Date',
  Deleted = 'Deleted'
}

export type QueryValue = {
  value?: InputMaybe<Scalars['Any']['input']>;
};

export enum RefRoot {
  AppRef = 'AppRef'
}

export enum RegisterAppCategory {
  App = 'App',
  Dashboard = 'Dashboard',
  Setting = 'Setting'
}

export type RegisterAppInput = {
  apiUrl: Scalars['String']['input'];
  apiWsUrl: Scalars['String']['input'];
  name: Scalars['String']['input'];
  navigation?: InputMaybe<Array<InputMaybe<RegisterAppNavigationInput>>>;
  package: Scalars['String']['input'];
  proxyApi: Scalars['Boolean']['input'];
  version: Scalars['String']['input'];
};

export type RegisterAppModule = {
  exposedModule?: InputMaybe<Scalars['String']['input']>;
  moduleName?: InputMaybe<Scalars['String']['input']>;
  outlet?: InputMaybe<Scalars['String']['input']>;
  path?: InputMaybe<Scalars['String']['input']>;
};

export type RegisterAppNavigationInput = {
  authRequired?: InputMaybe<Scalars['Boolean']['input']>;
  category: RegisterAppCategory;
  children?: InputMaybe<Array<InputMaybe<RegisterChildAppNavigationInput>>>;
  hidden?: InputMaybe<Scalars['Boolean']['input']>;
  icon: Scalars['String']['input'];
  module: RegisterAppModule;
  proxy: Scalars['Boolean']['input'];
  subTitle?: InputMaybe<Scalars['String']['input']>;
  title: Scalars['String']['input'];
};

export type RegisterAppOutput = {
  __typename?: 'RegisterAppOutput';
  id: Scalars['String']['output'];
};

export type RegisterAppSlot = {
  authRequired?: InputMaybe<Scalars['Boolean']['input']>;
  description: Scalars['String']['input'];
  module: RegisterAppSlotModule;
};

export type RegisterAppSlotModule = {
  exposedModule?: InputMaybe<Scalars['String']['input']>;
  moduleName?: InputMaybe<Scalars['String']['input']>;
  path?: InputMaybe<Scalars['String']['input']>;
};

export type RegisterChildAppNavigationInput = {
  authRequired?: InputMaybe<Scalars['Boolean']['input']>;
  children?: InputMaybe<Array<InputMaybe<RegisterChildAppNavigationInput>>>;
  icon: Scalars['String']['input'];
  module: RegisterAppModule;
  path?: InputMaybe<Scalars['String']['input']>;
  subTitle?: InputMaybe<Scalars['String']['input']>;
  title: Scalars['String']['input'];
};

export type RegisteredApp = {
  __typename?: 'RegisteredApp';
  createdAt?: Maybe<Scalars['Time']['output']>;
  name?: Maybe<Scalars['String']['output']>;
  pkg: Scalars['String']['output'];
  updatedAt?: Maybe<Scalars['Time']['output']>;
};

export type RegisteredAppQueryFields = {
  createdAt?: InputMaybe<QueryOperatorAndDate>;
  name?: InputMaybe<QueryOperatorAndValue>;
  updatedAt?: InputMaybe<QueryOperatorAndDate>;
};

export type RegisteredAppsPage = {
  __typename?: 'RegisteredAppsPage';
  data?: Maybe<Array<Maybe<RegisteredApp>>>;
  page: PageInfo;
};

export type RegisteredAppsSort = {
  createdAt?: InputMaybe<SortType>;
  name?: InputMaybe<SortType>;
  updatedAt?: InputMaybe<SortType>;
};

export type RegisteredAppsWhereRules = {
  condition: QueryCondition;
  fields?: InputMaybe<Array<InputMaybe<RegisteredAppQueryFields>>>;
  rules?: InputMaybe<Array<InputMaybe<RegisteredAppsWhereRules>>>;
};

export enum ShellConfigEventType {
  Added = 'Added',
  Initial = 'Initial',
  Rebuild = 'Rebuild',
  Removed = 'Removed',
  Updated = 'Updated'
}

export type ShellConfiguration = {
  __typename?: 'ShellConfiguration';
  categories?: Maybe<Array<Maybe<ShellNavigationCategory>>>;
  defaultRoute?: Maybe<Scalars['String']['output']>;
  slots?: Maybe<Array<Maybe<ShellNavigationSlot>>>;
};

export type ShellConfigurationSubscription = {
  __typename?: 'ShellConfigurationSubscription';
  configuration: ShellConfiguration;
  eventType: ShellConfigEventType;
};

export type ShellNavigation = {
  __typename?: 'ShellNavigation';
  authRequired?: Maybe<Scalars['Boolean']['output']>;
  children?: Maybe<Array<Maybe<ShellNavigationChild>>>;
  healthy: Scalars['Boolean']['output'];
  hidden: Scalars['Boolean']['output'];
  icon: Scalars['String']['output'];
  id: Scalars['String']['output'];
  module: ShellNavigationModule;
  subTitle?: Maybe<Scalars['String']['output']>;
  title: Scalars['String']['output'];
};

export type ShellNavigationCategory = {
  __typename?: 'ShellNavigationCategory';
  category: RegisterAppCategory;
  entries?: Maybe<Array<Maybe<ShellNavigation>>>;
  priority: Scalars['Int']['output'];
  title: Scalars['String']['output'];
};

export type ShellNavigationChild = {
  __typename?: 'ShellNavigationChild';
  authRequired?: Maybe<Scalars['Boolean']['output']>;
  children?: Maybe<Array<Maybe<ShellNavigationChild>>>;
  healthy: Scalars['Boolean']['output'];
  icon: Scalars['String']['output'];
  module: ShellNavigationModule;
  subTitle?: Maybe<Scalars['String']['output']>;
  title: Scalars['String']['output'];
};

export type ShellNavigationModule = {
  __typename?: 'ShellNavigationModule';
  exposedModule: Scalars['String']['output'];
  moduleName: Scalars['String']['output'];
  outlet: Scalars['String']['output'];
  path: Scalars['String']['output'];
  remoteEntry: Scalars['String']['output'];
};

export type ShellNavigationSlot = {
  __typename?: 'ShellNavigationSlot';
  authRequired?: Maybe<Scalars['Boolean']['output']>;
  description: Scalars['String']['output'];
  module: ShellNavigationSlotModule;
  priority?: Maybe<Scalars['Int']['output']>;
  slot: Scalars['String']['output'];
};

export type ShellNavigationSlotModule = {
  __typename?: 'ShellNavigationSlotModule';
  exposedModule: Scalars['String']['output'];
  moduleName: Scalars['String']['output'];
  path: Scalars['String']['output'];
  remoteEntry: Scalars['String']['output'];
};

export type Sort = {
  key: Scalars['String']['input'];
  type: SortType;
};

export enum SortType {
  Asc = 'ASC',
  Des = 'DES'
}

export type Subscription = {
  __typename?: 'Subscription';
  shellConfiguration: ShellConfigurationSubscription;
};


export type SubscriptionShellConfigurationArgs = {
  events: Array<ShellConfigEventType>;
  tenantId: Scalars['String']['input'];
};

export type TagValue = {
  __typename?: 'TagValue';
  Value: Scalars['Any']['output'];
};

export type TagValues = {
  __typename?: 'TagValues';
  Key: Scalars['String']['output'];
  Values?: Maybe<Array<Maybe<TagValue>>>;
};

export type FetchShellConfigQueryVariables = Exact<{
  tenant: Scalars['String']['input'];
}>;


export type FetchShellConfigQuery = { __typename?: 'Query', shellConfiguration: { __typename?: 'ShellConfiguration', defaultRoute?: string | null, categories?: Array<{ __typename?: 'ShellNavigationCategory', category: RegisterAppCategory, priority: number, title: string, entries?: Array<{ __typename?: 'ShellNavigation', id: string, title: string, subTitle?: string | null, authRequired?: boolean | null, healthy: boolean, hidden: boolean, icon: string, module: { __typename?: 'ShellNavigationModule', exposedModule: string, moduleName: string, outlet: string, path: string, remoteEntry: string } } | null> | null } | null> | null, slots?: Array<{ __typename?: 'ShellNavigationSlot', authRequired?: boolean | null, priority?: number | null, description: string, slot: string, module: { __typename?: 'ShellNavigationSlotModule', remoteEntry: string, path: string, moduleName: string, exposedModule: string } } | null> | null } };

export type SubscribeToShellConfigSubscriptionVariables = Exact<{
  tenant: Scalars['String']['input'];
  events: Array<ShellConfigEventType> | ShellConfigEventType;
}>;


export type SubscribeToShellConfigSubscription = { __typename?: 'Subscription', shellConfiguration: { __typename?: 'ShellConfigurationSubscription', eventType: ShellConfigEventType, configuration: { __typename?: 'ShellConfiguration', defaultRoute?: string | null, categories?: Array<{ __typename?: 'ShellNavigationCategory', category: RegisterAppCategory, priority: number, title: string, entries?: Array<{ __typename?: 'ShellNavigation', id: string, title: string, subTitle?: string | null, authRequired?: boolean | null, healthy: boolean, hidden: boolean, icon: string, module: { __typename?: 'ShellNavigationModule', exposedModule: string, moduleName: string, outlet: string, path: string, remoteEntry: string } } | null> | null } | null> | null, slots?: Array<{ __typename?: 'ShellNavigationSlot', authRequired?: boolean | null, priority?: number | null, description: string, slot: string, module: { __typename?: 'ShellNavigationSlotModule', remoteEntry: string, path: string, moduleName: string, exposedModule: string } } | null> | null } } };


export const FetchShellConfigDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"FetchShellConfig"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"tenant"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"shellConfiguration"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"tenantId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"tenant"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"defaultRoute"}},{"kind":"Field","name":{"kind":"Name","value":"categories"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"category"}},{"kind":"Field","name":{"kind":"Name","value":"priority"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"entries"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"subTitle"}},{"kind":"Field","name":{"kind":"Name","value":"authRequired"}},{"kind":"Field","name":{"kind":"Name","value":"healthy"}},{"kind":"Field","name":{"kind":"Name","value":"hidden"}},{"kind":"Field","name":{"kind":"Name","value":"icon"}},{"kind":"Field","name":{"kind":"Name","value":"module"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"exposedModule"}},{"kind":"Field","name":{"kind":"Name","value":"moduleName"}},{"kind":"Field","name":{"kind":"Name","value":"outlet"}},{"kind":"Field","name":{"kind":"Name","value":"path"}},{"kind":"Field","name":{"kind":"Name","value":"remoteEntry"}}]}}]}}]}},{"kind":"Field","name":{"kind":"Name","value":"slots"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"authRequired"}},{"kind":"Field","name":{"kind":"Name","value":"priority"}},{"kind":"Field","name":{"kind":"Name","value":"description"}},{"kind":"Field","name":{"kind":"Name","value":"slot"}},{"kind":"Field","name":{"kind":"Name","value":"module"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"remoteEntry"}},{"kind":"Field","name":{"kind":"Name","value":"path"}},{"kind":"Field","name":{"kind":"Name","value":"moduleName"}},{"kind":"Field","name":{"kind":"Name","value":"exposedModule"}}]}}]}}]}}]}}]} as unknown as DocumentNode<FetchShellConfigQuery, FetchShellConfigQueryVariables>;
export const SubscribeToShellConfigDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"subscription","name":{"kind":"Name","value":"SubscribeToShellConfig"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"tenant"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"events"}},"type":{"kind":"NonNullType","type":{"kind":"ListType","type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ShellConfigEventType"}}}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"shellConfiguration"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"tenantId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"tenant"}}},{"kind":"Argument","name":{"kind":"Name","value":"events"},"value":{"kind":"Variable","name":{"kind":"Name","value":"events"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"configuration"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"defaultRoute"}},{"kind":"Field","name":{"kind":"Name","value":"categories"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"category"}},{"kind":"Field","name":{"kind":"Name","value":"priority"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"entries"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"subTitle"}},{"kind":"Field","name":{"kind":"Name","value":"authRequired"}},{"kind":"Field","name":{"kind":"Name","value":"healthy"}},{"kind":"Field","name":{"kind":"Name","value":"hidden"}},{"kind":"Field","name":{"kind":"Name","value":"icon"}},{"kind":"Field","name":{"kind":"Name","value":"module"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"exposedModule"}},{"kind":"Field","name":{"kind":"Name","value":"moduleName"}},{"kind":"Field","name":{"kind":"Name","value":"outlet"}},{"kind":"Field","name":{"kind":"Name","value":"path"}},{"kind":"Field","name":{"kind":"Name","value":"remoteEntry"}}]}}]}}]}},{"kind":"Field","name":{"kind":"Name","value":"slots"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"authRequired"}},{"kind":"Field","name":{"kind":"Name","value":"priority"}},{"kind":"Field","name":{"kind":"Name","value":"description"}},{"kind":"Field","name":{"kind":"Name","value":"slot"}},{"kind":"Field","name":{"kind":"Name","value":"module"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"remoteEntry"}},{"kind":"Field","name":{"kind":"Name","value":"path"}},{"kind":"Field","name":{"kind":"Name","value":"moduleName"}},{"kind":"Field","name":{"kind":"Name","value":"exposedModule"}}]}}]}}]}},{"kind":"Field","name":{"kind":"Name","value":"eventType"}}]}}]}}]} as unknown as DocumentNode<SubscribeToShellConfigSubscription, SubscribeToShellConfigSubscriptionVariables>;