import {FC, useEffect, useState} from "react";

type VersionResponse = {
  service: string
  version: string
  build: string
}

export const HomePage: FC = () => {
  const [version, setVersion] = useState("")

  useEffect(() => {
    const fetchVersion = async () => {
      const resp = await fetch('/api/version');
      if (!resp.ok) {
        throw resp.statusText
      }
      const versionResp: VersionResponse = await resp.json() as VersionResponse
      setVersion(versionResp.service + " " + versionResp.version + " (" + versionResp.build + ")");
    };

    fetchVersion()
      .catch((e) => {
        console.log("failed to fetch version:", e.toString())
        setVersion("failed to fetch version from the server: " + e.toString())
      })
  }, [])

  return (
    <div id="workspace">

      <div className="header">
        HOME PAGE HEADER
      </div>

      <div className="content">
        API Service version: <b>{version}</b>
      </div>

    </div>
  );
}
