// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

import React, {useMemo} from "react"

import {Badge, Panel} from "@cloudoperators/juno-ui-components"
import EditDomainPanel from "./Domains/EditDomainPanel"
import EditDatacenterPanel from "./Datacenters/EditDatacenterPanel";
import EditPoolPanel from "./Pools/EditPoolPanel";
import EditMemberPanel from "./Pools/Members/EditMemberPanel";
import EditMonitorPanel from "./Pools/Monitors/EditMonitorPanel";
import {urlStore} from "../store";

const PanelManager = ({ currentPanel, closePanel }) => {
  const id = urlStore((state) => state.id)

  const heading = useMemo(() => {
    return (
        <span>
          Edit {currentPanel}: <Badge>{id}</Badge>
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
