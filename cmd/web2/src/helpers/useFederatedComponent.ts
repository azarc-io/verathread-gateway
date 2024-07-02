import loadComponent from './componentLoaderFromWebpackContainer';
import { getOrLoadRemote } from './getOrLoadRemote';
import { IFederatedComponent } from './types';
import {defineAsyncComponent} from "vue";

export default function useFederatedComponent({ remoteName, remoteUrl, moduleToLoad, shareScope = 'default' }: IFederatedComponent) {
  const key = `${remoteUrl}-${remoteName}-${shareScope}-${moduleToLoad}`;

  return () => new Promise(async (resolve, reject) => {
    await getOrLoadRemote({ remoteName, remoteUrl, shareScope });
    const com = loadComponent({ remoteName, moduleToLoad })
    resolve(com())
  })
}
