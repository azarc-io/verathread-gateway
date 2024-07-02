
    export type RemoteKeys = 'example/Counter';
    type PackageType<T> = T extends 'example/Counter' ? typeof import('example/Counter') :any;