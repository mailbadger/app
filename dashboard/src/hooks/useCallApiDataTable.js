import { useApi } from "./"

export const useCallApiDataTable = (endpoint) => {
    const [state, callApi] = useApi(
        {
            url: endpoint,
        },
        {
            collection: [],
            init: true,
        }
    )

    const onClickPrev = () => {
        callApi({
            url: state.data.links.previous,
        })
    }

    const onClickNext = () => {
        callApi({
            url: state.data.links.next,
        })
    }

    return [callApi, state, onClickPrev, onClickNext]
}
