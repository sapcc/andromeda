import React from "react"
import LogInModal from "./LogInModal"
import NewDomainModal from "./Domains/NewDomainModal"
import NewMemberModal from "./Pools/Members/NewMemberModal";
import NewMonitorModal from "./Pools/Monitors/NewMonitorModal";
import NewPoolModal from "./Pools/NewPoolModal";
import NewDatacenterModal from "./Datacenters/NewDatacenterModal";
import NewGeographicMapModal from "./GeographicMaps/NewGeographicMapModal";
import ShowGeographicMapModal from "./GeographicMaps/ShowGeographicMapModal";
import {urlStore} from "../store";

const ModalManager = () => {
    const currentModal = urlStore((state) => state.m)
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
            return <ShowGeographicMapModal />
        case "NewMembersItem":
            return <NewMemberModal/>
        case "NewMonitorsItem":
            return <NewMonitorModal/>
        default:
            return null
    }
}

export default ModalManager
