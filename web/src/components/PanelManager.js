import React, {useMemo} from "react"

import {Badge, Panel} from "juno-ui-components"
import {useStore} from "../store"
import {currentState} from "url-state-provider"
import EditDomainPanel from "./Domains/EditDomainPanel"
import EditDatacenterPanel from "./Datacenters/EditDatacenterPanel";
import EditPoolPanel from "./Pools/EditPoolPanel";
import EditMemberPanel from "./Pools/Members/EditMemberPanel";
import EditMonitorPanel from "./Pools/Monitors/EditMonitorPanel";

const PanelManager = ({ currentPanel, closePanel }) => {
  const urlStateKey = useStore((state) => state.urlStateKey)

  const heading = useMemo(() => {
    const urlState = currentState(urlStateKey)
    return (
        <span>
          Edit {currentPanel}: <Badge>{urlState?.id}</Badge>
        </span>
    )
  }, [currentPanel])

  const panelBody = () => {
    switch (currentPanel) {
      case "Domain":
        return <EditDomainPanel closeCallback={closePanel} />
      case "Datacenter":
        return <EditDatacenterPanel closeCallback={closePanel} />
      case "Pool":
        return <EditPoolPanel closeCallback={closePanel} />
      case "Member":
        return <EditMemberPanel closeCallback={closePanel} />
      case "Monitor":
        return <EditMonitorPanel closeCallback={closePanel} />
      default:
        return null
    }
  }

  return (
    <Panel
      heading={heading}
      opened={!!panelBody()}
      onClose={closePanel}
    >
      {panelBody()}
    </Panel>
  )
}

export default PanelManager
