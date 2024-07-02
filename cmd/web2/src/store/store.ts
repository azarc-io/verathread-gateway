import {reactive} from "vue";
import {FetchShellConfigDocument} from "../gql/graphql";

export interface State {
    configuration: FetchShellConfigDocument
}

const state = reactive({
    configuration: FetchShellConfigDocument,
})

export default state
