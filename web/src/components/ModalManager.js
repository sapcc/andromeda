import React from "react"
import LogInModal from "./LogInModal"
import NewDomainModal from "./Domains/NewDomainModal"
import NewMemberModal from "./Pools/Members/NewMemberModal";
import NewMonitorModal from "./Pools/Monitors/NewMonitorModal";
import NewPoolModal from "./Pools/NewPoolModal";
import NewDatacenterModal from "./Datacenters/NewDatacenterModal";

const ModalManager = ({currentModal}) => {
    switch (currentModal) {
        case "LogIn":
            return <LogInModal />
        case "NewDomainsItem":
            return <NewDomainModal />
        case "NewPoolsItem":
            return <NewPoolModal />
        case "NewDatacentersItem":
            return <NewDatacenterModal />
        case "NewMembersItem":
            return <NewMemberModal/>
        case "NewMonitorsItem":
            return <NewMonitorModal/>
        default:
            return null
    }
}

export default ModalManager
