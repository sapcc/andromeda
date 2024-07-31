import React from "react"

import {Container, IntroBox, MainTabs, Tab, TabList, TabPanel,} from "@cloudoperators/juno-ui-components"
import {urlStore} from "./store"
import ModalManager from "./components/ModalManager"
import PanelManager from "./components/PanelManager"
import DomainList from "./components/Domains/DomainList"
import DatacenterList from "./components/Datacenters/DatacenterList";
import PoolList from "./components/Pools/PoolList";
import GeographicMapList from "./components/GeographicMaps/GeographicMapList";

const AppContent = ({props}) => {
  const [panel, setPanel] = urlStore((state) => [state.p, state.openPanel])
  const [tab, setTab] = urlStore((state) => [state.t, state.setTab])

  return (
    <>
      <MainTabs selectedIndex={tab} onSelect={setTab}>
        <PanelManager currentPanel={panel} closePanel={() => setPanel(null)} />

        <TabList>
          <Tab>Domains</Tab>
          <Tab>Pools</Tab>
          <Tab>Datacenters</Tab>
          <Tab>Geographic Maps</Tab>
        </TabList>

        <TabPanel>
          {/* You'll normally want to use a Container as a wrapper for your content because it has padding that makes everything look nice */}
          <Container py>
            <IntroBox
                title="Andromeda"
                heroImage="url(data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNDgzIiBoZWlnaHQ9IjIzMSIgdmlld0JveD0iMCAwIDQ4MyAyMzEiIGZpbGw9InJnYmEoMzYsIDQyLCA0OSwgMSkiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CjxtYXNrIGlkPSJtYXNrMF8zOTZfMjA0ODUiIHN0eWxlPSJtYXNrLXR5cGU6YWxwaGEiIG1hc2tVbml0cz0idXNlclNwYWNlT25Vc2UiIHg9IjAiIHk9IjAiIHdpZHRoPSI0ODMiIGhlaWdodD0iMjMxIj4KPHJlY3Qgd2lkdGg9IjQ4MyIgaGVpZ2h0PSIyMzEiIGZpbGw9IiNDNEM0QzQiLz4KPC9tYXNrPgo8ZyBtYXNrPSJ1cmwoI21hc2swXzM5Nl8yMDQ4NSkiPgo8cGF0aCBmaWxsLXJ1bGU9ImV2ZW5vZGQiIGNsaXAtcnVsZT0iZXZlbm9kZCIgZD0iTTM1NC43NDMgLTExLjE0NDRDMzU0LjcwNyAtMTEuMjY3OSAzNTQuNjc0IC0xMS4zOTM0IDM1NC42NDYgLTExLjUyMDlDMzUzLjg1NyAtMTUuMDI2NCAzNTYuMjQ5IC0xOC43MDIgMzU5Ljk4NyAtMTkuNzMwNkMzNjMuNzI0IC0yMC43NTkzIDM2Ny4zOTQgLTE4Ljc1MTUgMzY4LjE4MiAtMTUuMjQ2MUMzNjguOTcgLTExLjc0MDcgMzY2LjU3OSAtOC4wNjUwNSAzNjIuODQxIC03LjAzNjM5QzM2MC43NTEgLTYuNDYxMDUgMzU4LjY4MSAtNi44MzU1OCAzNTcuMTI2IC03Ljg5MDc0TDI3MS4yMyAzMC4yMTkyTDMxMS4zNzcgNDMuMzA1NkMzMTEuNzk1IDQ0LjM4ODMgMzEyLjczMSA0NS4wNzIgMzEzLjQwOCA0NS40MjMxQzMxNC40MTkgNDUuOTQ3NyAzMTUuNjU1IDQ2LjIwMzEgMzE2Ljg5NCA0Ni4yMDMxQzMxNy43MDIgNDYuMjAzMSAzMTguNTA5IDQ2LjA5NDQgMzE5LjI1NCA0NS44NzMyTDMyMC4xNDMgNDYuMTYzTDMzOC41MTggMzcuODExNFY0Mi4xMzk4TDMyNS42MTggNDguMDAyOEwzMzMuNDIxIDUwLjYzNjhMMzI3Ljk2MiA1Mi45NTIxTDMxOS42NDggNTAuMTQ1NEwzMTkuNjI4IDUwLjEzODlMMjY5LjYyOCAzMy44NDA3QzI2Ni4zMjMgMzIuNzYzNCAyNjYuMDg3IDI4LjE5IDI2OS4yNjQgMjYuNzgwNUwzNTQuNzQzIC0xMS4xNDQ0Wk0zNjEuOTc3IC0xMC44Nzk2QzM2MC4wNzIgLTEwLjM1NTMgMzU4LjcwNCAtMTEuNDQ3OCAzNTguNDUyIC0xMi41Njg1QzM1OC4yIC0xMy42ODkxIDM1OC45NDYgLTE1LjM2MzEgMzYwLjg1MSAtMTUuODg3NEMzNjIuNzU2IC0xNi40MTE3IDM2NC4xMjMgLTE1LjMxOTIgMzY0LjM3NSAtMTQuMTk4NUMzNjQuNjI3IC0xMy4wNzc5IDM2My44ODIgLTExLjQwMzkgMzYxLjk3NyAtMTAuODc5NloiLz4KPHBhdGggZmlsbC1ydWxlPSJldmVub2RkIiBjbGlwLXJ1bGU9ImV2ZW5vZGQiIGQ9Ik0zNDkuOTg4IDM2LjkyNjRWMzIuNTk4TDQ0NS4wNjUgLTEwLjYxNTNDNDQ0LjI3OSAtMTQuMTE5OCA0NDYuNjY5IC0xNy43OTM0IDQ1MC40MDYgLTE4LjgyMTdDNDU0LjE0NCAtMTkuODUwNCA0NTcuODEzIC0xNy44NDI2IDQ1OC42MDEgLTE0LjMzNzFDNDU5LjM5IC0xMC44MzE3IDQ1Ni45OTkgLTcuMTU2MTIgNDUzLjI2MSAtNi4xMjc0NUM0NTEuMDAyIC01LjUwNTkzIDQ0OC43NjkgLTUuOTkyODkgNDQ3LjE4MSAtNy4yNDg5MUwzNDkuOTg4IDM2LjkyNjRaTTQ0OC44NzIgLTExLjY1OTZDNDQ5LjEyNCAtMTAuNTM4OSA0NTAuNDkxIC05LjQ0NjQgNDUyLjM5NiAtOS45NzA2OEM0NTQuMzAyIC0xMC40OTUgNDU1LjA0NyAtMTIuMTY4OSA0NTQuNzk1IC0xMy4yODk2QzQ1NC41NDMgLTE0LjQxMDIgNDUzLjE3NSAtMTUuNTAyNyA0NTEuMjcgLTE0Ljk3ODVDNDQ5LjM2NSAtMTQuNDU0MiA0NDguNjIgLTEyLjc4MDIgNDQ4Ljg3MiAtMTEuNjU5NloiLz4KPHBhdGggZmlsbC1ydWxlPSJldmVub2RkIiBjbGlwLXJ1bGU9ImV2ZW5vZGQiIGQ9Ik0zNjQuMjY1IDY1LjIwNzZWNjEuMDQ5NEwzNzUuOTgyIDY1LjAwNDdMNDI0Ljc5NiA0NC4xMzcyQzQyNC4wNTYgNDAuNjUzNCA0MjYuNDQgMzcuMDIwNyA0MzAuMTUzIDM1Ljk5OUM0MzMuODkxIDM0Ljk3MDMgNDM3LjU2IDM2Ljk3ODEgNDM4LjM0OCA0MC40ODM2QzQzOS4xMzcgNDMuOTg5IDQzNi43NDUgNDcuNjY0NiA0MzMuMDA4IDQ4LjY5MzJDNDMwLjcyNiA0OS4zMjEyIDQyOC40NjkgNDguODE3NiA0MjYuODc4IDQ3LjUzMjRMNDA2LjM5OSA1Ni4yODdMNDMwLjcxNCA2My44NTFMNDg0Ljg1OSA0MC4zODE1QzQ4NC43NzYgNDAuMTUwMSA0ODQuNzA2IDM5LjkxMDggNDg0LjY1MSAzOS42NjRDNDgzLjg2MyAzNi4xNTg2IDQ4Ni4yNTQgMzIuNDgzIDQ4OS45OTIgMzEuNDU0M0M0OTMuNzMgMzAuNDI1NiA0OTcuMzk5IDMyLjQzMzUgNDk4LjE4NyAzNS45Mzg5QzQ5OC45NzUgMzkuNDQ0MyA0OTYuNTg0IDQzLjExOTkgNDkyLjg0NiA0NC4xNDg2QzQ5MC45MTkgNDQuNjc4OSA0ODkuMDEgNDQuNDAyMSA0ODcuNTA1IDQzLjUyOTJMNDMyLjIzNSA2Ny40ODU5QzQzMS40IDY3Ljg0ODEgNDMwLjQ2MiA2Ny44OTg1IDQyOS41OTIgNjcuNjI3OUw0MDAuODEzIDU4LjY3NTFMMzgxLjU4MyA2Ni44OTU1TDQ0OS4yODggODkuNzUxOUw0OTYuMzUxIDY4LjcwOTFDNDk2LjMzOCA2OC42NTc3IDQ5Ni4zMjUgNjguNjA1OSA0OTYuMzE0IDY4LjU1MzhDNDk1LjUyNSA2NS4wNDg0IDQ5Ny45MTYgNjEuMzcyOCA1MDEuNjU0IDYwLjM0NDFDNTA1LjM5MiA1OS4zMTU1IDUwOS4wNjEgNjEuMzIzMyA1MDkuODUgNjQuODI4N0M1MTAuNjM4IDY4LjMzNDEgNTA4LjI0NyA3Mi4wMDk3IDUwNC41MDkgNzMuMDM4NEM1MDIuMzE5IDczLjY0MTEgNTAwLjE1MiA3My4yMDE0IDQ5OC41NzUgNzIuMDI4OEw0NTAuNzM5IDkzLjQxNTVDNDQ5Ljg5MSA5My43NjM5IDQ0OC45NDUgOTMuNzk0MyA0NDguMDc2IDkzLjUwMUwzNzYuMTQ4IDY5LjIxODlMMzY0LjI2NSA3My44NDg3VjY5LjY1OTJMMzcwLjU0NyA2Ny4zMjgxTDM2NC4yNjUgNjUuMjA3NlpNNDMyLjE0MyA0NC44NUM0MzAuMjM4IDQ1LjM3NDMgNDI4Ljg3MSA0NC4yODE4IDQyOC42MTkgNDMuMTYxMUM0MjguMzY3IDQyLjA0MDUgNDI5LjExMiA0MC4zNjY1IDQzMS4wMTcgMzkuODQyMkM0MzIuOTIyIDM5LjMxNzkgNDM0LjI5IDQwLjQxMDUgNDM0LjU0MiA0MS41MzExQzQzNC43OTQgNDIuNjUxOCA0MzQuMDQ4IDQ0LjMyNTcgNDMyLjE0MyA0NC44NVpNNDkxLjk4MiA0MC4zMDUzQzQ5MC4wNzcgNDAuODI5NiA0ODguNzA5IDM5LjczNzEgNDg4LjQ1NyAzOC42MTY1QzQ4OC4yMDUgMzcuNDk1OCA0ODguOTUxIDM1LjgyMTggNDkwLjg1NiAzNS4yOTc1QzQ5Mi43NjEgMzQuNzczMyA0OTQuMTI4IDM1Ljg2NTggNDk0LjM4MSAzNi45ODY0QzQ5NC42MzMgMzguMTA3MSA0OTMuODg3IDM5Ljc4MTEgNDkxLjk4MiA0MC4zMDUzWk01MDAuMTIgNjcuNTA2M0M1MDAuMzcyIDY4LjYyNjkgNTAxLjc0IDY5LjcxOTQgNTAzLjY0NSA2OS4xOTUxQzUwNS41NSA2OC42NzA5IDUwNi4yOTUgNjYuOTk2OSA1MDYuMDQzIDY1Ljg3NjJDNTA1Ljc5MSA2NC43NTU2IDUwNC40MjQgNjMuNjYzMSA1MDIuNTE5IDY0LjE4NzRDNTAwLjYxMyA2NC43MTE2IDQ5OS44NjggNjYuMzg1NiA1MDAuMTIgNjcuNTA2M1oiLz4KPHBhdGggZmlsbC1ydWxlPSJldmVub2RkIiBjbGlwLXJ1bGU9ImV2ZW5vZGQiIGQ9Ik0zMDEuMzY2IDkzLjAwMlY5OC4zNTY1TDI1OC42ODEgMTE0Ljk4OEwzNzguOTYxIDE1NS41OTNDMzgyLjA5MiAxNTYuNjUgMzgyLjQ2NyAxNjAuOTE2IDM3OS41NjkgMTYyLjUwMUwyODAuODA0IDIxNi41MDlDMjgwLjkzOSAyMTYuODMgMjgxLjA0NyAyMTcuMTY4IDI4MS4xMjYgMjE3LjUyMUMyODEuOTE1IDIyMS4wMjYgMjc5LjUyNCAyMjQuNzAyIDI3NS43ODYgMjI1LjczMUMyNzIuMDQ4IDIyNi43NTkgMjY4LjM3OSAyMjQuNzUxIDI2Ny41OSAyMjEuMjQ2QzI2Ni44MDIgMjE3Ljc0MSAyNjkuMTkzIDIxNC4wNjUgMjcyLjkzMSAyMTMuMDM2QzI3NC43NjEgMjEyLjUzMyAyNzYuNTc0IDIxMi43NTcgMjc4LjA0MSAyMTMuNTI4TDM3Ny4zNjcgMTU5LjIxM0wzMTcuODQxIDEzOS4xMThMMTcyLjA1NiAyMDkuMjM3QzE3Mi4wNiAyMDkuMjU2IDE3Mi4wNjUgMjA5LjI3NCAxNzIuMDY5IDIwOS4yOTNMMTcyLjA3NSAyMDkuMzE3TDE3Mi4wOCAyMDkuMzRDMTcyLjg2OCAyMTIuODQ2IDE3MC40NzcgMjE2LjUyMSAxNjYuNzM5IDIxNy41NUMxNjMuMDAxIDIxOC41NzkgMTU5LjMzMiAyMTYuNTcxIDE1OC41NDQgMjEzLjA2NkMxNTcuNzU2IDIwOS41NiAxNjAuMTQ3IDIwNS44ODUgMTYzLjg4NSAyMDQuODU2QzE2Ni4xMDMgMjA0LjI0NSAxNjguMjk3IDIwNC43MDQgMTY5Ljg3OSAyMDUuOTExTDMxMi40OTkgMTM3LjMxNUwyNTMuMjQ2IDExNy4zMTJMMTIyLjEyMiAxNzMuMzY1QzEyMi4xNzcgMTczLjUzNyAxMjIuMjI1IDE3My43MTIgMTIyLjI2NSAxNzMuODkyQzEyMy4wNTQgMTc3LjM5NyAxMjAuNjYzIDE4MS4wNzMgMTE2LjkyNSAxODIuMTAyQzExMy4xODcgMTgzLjEzIDEwOS41MTcgMTgxLjEyMiAxMDguNzI5IDE3Ny42MTdDMTA3Ljk0MSAxNzQuMTEyIDExMC4zMzIgMTcwLjQzNiAxMTQuMDcgMTY5LjQwN0MxMTYuMTA3IDE2OC44NDcgMTE4LjEyNCAxNjkuMTg4IDExOS42NjUgMTcwLjE4MkwxMTkuNjQ2IDE3MC4xMzhMMjg2LjM1MyA5OC44NzM2TDIzMy4yODcgNzkuODIyMkwxOC42NTAyIDE2NS41NTlDMTguNjU4OCAxNjUuNTk0IDE4LjY2NzIgMTY1LjYyOCAxOC42NzUzIDE2NS42NjNDMTguNjc5MSAxNjUuNjc5IDE4LjY4MjggMTY1LjY5NSAxOC42ODY0IDE2NS43MTFDMTkuNDc0NyAxNjkuMjE3IDE3LjA4MzYgMTcyLjg5MiAxMy4zNDU3IDE3My45MjFDOS42MDc3OCAxNzQuOTUgNS45Mzg1NyAxNzIuOTQyIDUuMTUwMjggMTY5LjQzN0M0LjM2MTk4IDE2NS45MzEgNi43NTMxMSAxNjIuMjU2IDEwLjQ5MSAxNjEuMjI3QzEyLjY2ODggMTYwLjYyOCAxNC44MjMyIDE2MS4wNTkgMTYuMzk4MiAxNjIuMjE3TDIzMS45MzIgNzYuMTIyMkMyMzIuNzggNzUuNzk5MyAyMzMuNzE3IDc1Ljc5MDQgMjM0LjU3MSA3Ni4wOTcyTDI5MS42NzYgOTYuNTk4NEwzMDEuMzY2IDkzLjAwMlpNMjc0LjkyMSAyMjEuODg3QzI3My4wMTYgMjIyLjQxMiAyNzEuNjQ5IDIyMS4zMTkgMjcxLjM5NyAyMjAuMTk5QzI3MS4xNDUgMjE5LjA3OCAyNzEuODkgMjE3LjQwNCAyNzMuNzk1IDIxNi44OEMyNzUuNyAyMTYuMzU1IDI3Ny4wNjggMjE3LjQ0OCAyNzcuMzIgMjE4LjU2OEMyNzcuNTcyIDIxOS42ODkgMjc2LjgyNyAyMjEuMzYzIDI3NC45MjEgMjIxLjg4N1pNMTYyLjM1IDIxMi4wMThDMTYyLjYwMiAyMTMuMTM5IDE2My45NyAyMTQuMjMxIDE2NS44NzUgMjEzLjcwN0MxNjcuNzggMjEzLjE4MyAxNjguNTI2IDIxMS41MDkgMTY4LjI3NCAyMTAuMzg4QzE2OC4wMjIgMjA5LjI2NyAxNjYuNjU0IDIwOC4xNzUgMTY0Ljc0OSAyMDguNjk5QzE2Mi44NDQgMjA5LjIyMyAxNjIuMDk4IDIxMC44OTcgMTYyLjM1IDIxMi4wMThaTTExNi4wNiAxNzguMjU4QzExNC4xNTUgMTc4Ljc4MyAxMTIuNzg4IDE3Ny42OSAxMTIuNTM2IDE3Ni41NjlDMTEyLjI4NCAxNzUuNDQ5IDExMy4wMjkgMTczLjc3NSAxMTQuOTM0IDE3My4yNTFDMTE2LjgzOSAxNzIuNzI2IDExOC4yMDcgMTczLjgxOSAxMTguNDU5IDE3NC45MzlDMTE4LjcxMSAxNzYuMDYgMTE3Ljk2NSAxNzcuNzM0IDExNi4wNiAxNzguMjU4Wk04Ljk1NjgxIDE2OC4zODlDOS4yMDg4MiAxNjkuNTEgMTAuNTc2MyAxNzAuNjAyIDEyLjQ4MTQgMTcwLjA3OEMxNC4zODY1IDE2OS41NTQgMTUuMTMxOSAxNjcuODggMTQuODc5OSAxNjYuNzU5QzE0LjYyNzkgMTY1LjYzOCAxMy4yNjA0IDE2NC41NDYgMTEuMzU1MyAxNjUuMDdDOS40NTAyIDE2NS41OTQgOC43MDQ4IDE2Ny4yNjggOC45NTY4MSAxNjguMzg5WiIvPgo8cGF0aCBmaWxsLXJ1bGU9ImV2ZW5vZGQiIGNsaXAtcnVsZT0iZXZlbm9kZCIgZD0iTTMyMi40MDUgMjAuMDMwMkMzMjAuODYyIDIyLjE2ODQgMzIwLjMxNCAyNC44MzY3IDMyMC4zMTQgMjYuNTY1VjQyLjExOTlDMzIwLjMxNCA0My4wOTk2IDMxOC43ODMgNDMuODkzNyAzMTYuODk0IDQzLjg5MzdDMzE1LjAwNSA0My44OTM3IDMxMy40NzQgNDMuMDk5NiAzMTMuNDc0IDQyLjExOTlWMjYuNTY1QzMxMy40NzQgMjMuNzQ1IDMxNC4yOTQgMTkuNTkxIDMxNi44NTUgMTYuMDQzOUMzMTkuNTQxIDEyLjMyMjkgMzIzLjk4NyA5LjUwOTE1IDMzMC41NzQgOS41MDkxNUMzMzcuMTYgOS41MDkxNSAzNDEuNjA2IDEyLjMyMjkgMzQ0LjI5MiAxNi4wNDM5QzM0Ni44NTMgMTkuNTkxIDM0Ny42NzMgMjMuNzQ1IDM0Ny42NzMgMjYuNTY1VjQ3LjIzMzdDMzQ3Ljc1NiA0Ny4yNDExIDM0Ny44MzkgNDcuMjU4OCAzNDcuOTE5IDQ3LjI4NjZMMzYxLjI1OSA1MS45MjAxQzM2MS40NjYgNTEuOTkyIDM2MS42MzkgNTIuMTI1NSAzNjEuNzYgNTIuMjk1N0wzNjEuOTUgNTIuMzYxNkwzNjEuODMyIDUyLjQxMTdDMzYxLjkwOCA1Mi41NTU4IDM2MS45NSA1Mi43MTkgMzYxLjk1IDUyLjg4OTRWMTAwLjM5OUMzNjEuOTUgMTAwLjgxMSAzNjEuNzAzIDEwMS4xODMgMzYxLjMyMyAxMDEuMzQzTDMxOC40MzkgMTE5LjQ5N0MzMTguMzI3IDExOS41NDUgMzE4LjIwOCAxMTkuNTcyIDMxOC4wODggMTE5LjU3OFYxMTkuODQ4TDMwMy42ODIgMTE1LjI2OVY2NS45NjE0TDMwMy44MzggNjUuODk1MkMzMDMuOTQ3IDY1LjcyMDIgMzA0LjExIDY1LjU3ODQgMzA0LjMwOCA2NS40OTQyTDMzOC42NzUgNTAuOTE3OEMzMzkuMTUgNTEuMDEyMyAzMzkuNjQzIDUxLjA1OSAzNDAuMTM3IDUxLjA1OUMzNDAuMzcgNTEuMDU5IDM0MC42MDMgNTEuMDQ4NiAzNDAuODMzIDUxLjAyNzhWMjYuNTY1QzM0MC44MzMgMjQuODM2NyAzNDAuMjg1IDIyLjE2ODQgMzM4Ljc0MiAyMC4wMzAyQzMzNy4zMjQgMTguMDY1OSAzMzQuOTMxIDE2LjMzMTUgMzMwLjU3NCAxNi4zMzE1QzMyNi4yMTYgMTYuMzMxNSAzMjMuODIzIDE4LjA2NTkgMzIyLjQwNSAyMC4wMzAyWk0zNjAuODAyIDUyLjg0ODNMMzYwLjkyMSA1Mi44ODk0VjEwMC4zOTlMMzE4LjA4OCAxMTguNTMxVjcwLjk2NTJMMzYwLjgwMiA1Mi44NDgzWk0zMjEuOTk0IDY1LjI4NkMzMjQuMTk2IDY1LjI4NiAzMjUuOTgxIDY0LjQyNDQgMzI1Ljk4MSA2My4zNjE1QzMyNS45ODEgNjIuMjk4NiAzMjQuMTk2IDYxLjQzNyAzMjEuOTk0IDYxLjQzN0MzMTkuNzkyIDYxLjQzNyAzMTguMDA2IDYyLjI5ODYgMzE4LjAwNiA2My4zNjE1QzMxOC4wMDYgNjQuNDI0NCAzMTkuNzkyIDY1LjI4NiAzMjEuOTk0IDY1LjI4NloiLz4KPC9nPgo8L3N2Zz4KCg==)"
                variant="hero"
            >
              <p>Platform agnostic GSLB frontend with OpenStack-like API</p>
              <p><small>Andromeda is in BETA! For support, please visit Slack Channel <a href={props.slackURL}>#andromeda</a> or send an <a href={`mailto:${props.helpMail}`}>email</a>.</small></p>
            </IntroBox>
            <DomainList />
          </Container>
        </TabPanel>
        <TabPanel>
          <Container py>
            <PoolList />

          </Container>
        </TabPanel>
        <TabPanel>
          <Container py>
            <DatacenterList />

          </Container>
        </TabPanel>
        <TabPanel>
          <Container py>
            <GeographicMapList />

          </Container>
        </TabPanel>
      </MainTabs>
      <ModalManager />
    </>
  )
}

export default AppContent
