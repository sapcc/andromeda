import React from "react"
import LogInModal from "./LogInModal"
import NewDomainModal from "./Domains/NewDomainModal"
import NewMemberModal from "./Pools/Members/NewMemberModal";
import NewMonitorModal from "./Pools/Monitors/NewMonitorModal";
import NewPoolModal from "./Pools/NewPoolModal";
import NewDatacenterModal from "./Datacenters/NewDatacenterModal";
import NewGeographicMapModal from "./GeographicMaps/NewGeographicMapModal";
import ShowGeographicMapModal from "./GeographicMaps/ShowGeographicMapModal";
import {useStore} from "../store";
import {currentState} from "url-state-provider";

const ModalManager = ({currentModal}) => {
    const urlStateKey = useStore((state) => state.urlStateKey)
    const urlState = currentState(urlStateKey)

    switch (currentModal) {
        case "LogIn":
            return <LogInModal />
        case "NewDomainsItem":
            return <NewDomainModal />
        case "NewPoolsItem":
            return <NewPoolModal />
        case "NewDatacentersItem":
            return <NewDatacenterModal />
        case "NewGeographicMapsItem":
            return <NewGeographicMapModal />
        case "ShowGeographicMap":
            return <ShowGeographicMapModal id={urlState?.id} />
        case "NewMembersItem":
            return <NewMemberModal/>
        case "NewMonitorsItem":
            return <NewMonitorModal/>
        default:
            return null
    }
}

export default ModalManager
