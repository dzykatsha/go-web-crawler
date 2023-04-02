import styles from "@/styles/Home.module.css";
import { Statistics, URLInfo } from "@/types";
import axios from "axios";
import Head from "next/head";
import { useRouter } from "next/router";
import { useEffect, useState } from "react";
import Modal from "react-modal";

const API_URL = process.env.NEXT_PUBLIC_API_URL;

export default function Home() {
  const [stats, setStats] = useState<Statistics>({ total: 0, urls: [] });
  const [page, setPage] = useState<number>(1);
  const [url, setUrl] = useState<string>();
  const [depth, setDepth] = useState<number>(0);
  const [urlInfo, setUrlInfo] = useState<URLInfo>();
  const [isOpen, setIsOpen] = useState<boolean>(false);

  const router = useRouter();

  const fetchStatistics = () => {
    axios
      .get(`${API_URL}/statistics`, {
        params: {
          page: page,
        },
      })
      .then((res) => setStats(res.data))
      .catch((err) => {
        setStats({ urls: [], total: 0 });
        alert(err);
      });
  };

  const sendUrl = async () => {
    try {
      const response = await axios.post(`${API_URL}/load`, {
        url: url,
        depth: depth,
      });
      alert(response.data);
    } catch (error) {
      alert(error);
    }
  };

  const redirectToPage = (uid: string) => {
    router.push({ pathname: "/page", query: { uid: uid } }, "/page");
  };

  const openHTML = async (uid: string) => {
    const result = await axios.get(`${API_URL}/page`, {
      params: { uid: uid },
    });
    setUrlInfo(result.data);
    setIsOpen(true);
  };

  useEffect(fetchStatistics, [page]);

  return (
    <>
      <Head>
        <title>Loaded pages</title>
        <meta name="description" content="Generated by create next app" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <main className={styles.main}>
        <div className="container my-5">
          <div className="container">
            <h1 className="p-2">Crawler statistics</h1>
            <div className="p-2">
              <a className="btn btn-primary" onMouseDown={fetchStatistics}>
                <i className="bi bi-arrow-clockwise"></i> Reload
              </a>
            </div>
          </div>
          <div className="container">
            <p className="fs-5">
              Totally downloaded{" "}
              <b>
                {stats && stats.total} {!stats && " Loading... "}
              </b>{" "}
              pages.
            </p>
            <table className="table">
              <thead>
                <tr>
                  <th scope="col">ID</th>
                  <th scope="col">Loaded time</th>
                  <th scope="col">URL</th>
                </tr>
              </thead>
              <tbody>
                {stats &&
                  stats.urls &&
                  stats.urls.map((url) => (
                    <tr key={url.uid}>
                      <th scope="row">
                        <a
                          className="link-info"
                          onMouseDown={() => openHTML(url.uid)}
                        >
                          {url.uid}
                        </a>
                      </th>
                      <td>
                        {new Date(url.createdAt * 1000).toLocaleTimeString(
                          "ru-RU"
                        )}
                      </td>
                      <td>{url.url}</td>
                    </tr>
                  ))}
              </tbody>
            </table>
          </div>
          <div className="container mt-2">
            <div className="row">
              {page != 1 && (
                <div className="col-md-auto">
                  <button
                    className="btn btn-primary"
                    onMouseDown={() => setPage(page - 1)}
                  >
                    Previous
                  </button>
                </div>
              )}
              <div className="col-md-auto">
                <button className="btn btn-light" disabled>
                  {page}
                </button>
              </div>
              {stats && stats.total > page * 10 && (
                <div className="col-md-auto">
                  <button
                    className="btn btn-primary"
                    onMouseDown={() => setPage(page + 1)}
                  >
                    Next
                  </button>
                </div>
              )}
            </div>
          </div>
          <div className="container-md mt-5">
            <h3>Download new page</h3>
            <form name="contact-form">
              <div className="row w-50">
                <div className="form-group">
                  <input
                    id="url"
                    name="url"
                    type="url"
                    className="form-control"
                    placeholder="Enter url of web page to scarp"
                    value={url}
                    onChange={(e) => setUrl(e.target.value)}
                    required
                  />
                </div>
                <div className="form-group xs-4">
                  <label className="form-label">Depth: {depth}</label>
                  <input
                    type="range"
                    id="depth"
                    name="depth"
                    className="form-range"
                    min="0"
                    max="5"
                    step="1"
                    value={depth}
                    onChange={(e) => setDepth(+e.target.value)}
                  ></input>
                </div>
              </div>
              <button
                type="submit"
                className="btn btn-primary"
                onMouseDown={sendUrl}
              >
                Submit
              </button>
            </form>
          </div>
        </div>
        <Modal
          isOpen={isOpen}
          onRequestClose={() => setIsOpen(false)}
          ariaHideApp={false}
        >
          <div className="container">
            <div className="row">
              <button
                className="col-md-auto btn btn-danger"
                onClick={() => setIsOpen(false)}
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="16"
                  height="16"
                  fill="currentColor"
                  className="bi bi-x-lg"
                  viewBox="0 0 16 16"
                >
                  <path d="M2.146 2.854a.5.5 0 1 1 .708-.708L8 7.293l5.146-5.147a.5.5 0 0 1 .708.708L8.707 8l5.147 5.146a.5.5 0 0 1-.708.708L8 8.707l-5.146 5.147a.5.5 0 0 1-.708-.708L7.293 8 2.146 2.854Z" />
                </svg>
              </button>
              <h4 className="col-md-auto text-truncate">
                {urlInfo && urlInfo.url}
              </h4>
            </div>
            <hr />
            <div className="row">
              <div className="overflow-auto h-100">
                <p className="text-monospace">{urlInfo?.html}</p>
              </div>
            </div>
          </div>
        </Modal>
      </main>
    </>
  );
}
