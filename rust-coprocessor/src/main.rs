use hyper::body::{self};
use hyper::service::{make_service_fn, service_fn};
use hyper::{Body, Request, Response, Server};
use json;
use std::convert::Infallible;
use std::net::SocketAddr;


#[tokio::main]
async fn main() {
    // We'll bind to 127.0.0.1:3000
    let addr = SocketAddr::from(([127, 0, 0, 1], 3007));

    let make_svc =
        make_service_fn(|_conn| async { Ok::<_, Infallible>(service_fn(request_handler)) });

    let server = Server::bind(&addr).serve(make_svc);

    println!("Listening on http://{}", addr);

    if let Err(e) = server.await {
        eprintln!("server error: {}", e);
    }
}

async fn request_handler(_req: Request<Body>) -> Result<Response<Body>, Infallible> {
    let body_contents = body::to_bytes(_req.into_body()).await.unwrap();
    let body = String::from_utf8(body_contents.to_vec()).unwrap();
    let parsed_body = json::parse(body.clone().as_str()).unwrap();
    
    return match parsed_body["stage"].as_str().unwrap() {
        "RouterRequest" => router_request(body),
        _ => Ok(Response::new(body.into())),
    };
}

fn router_request(
    body: String,
) -> Result<Response<Body>, Infallible> {
    println!("RouterRequest: Paylaod {:?}", body);
    return Ok(Response::new(body.into()));
}
