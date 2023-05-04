use polylib as poly;
use std::fs::File;
use std::os::unix::prelude::FileExt;
use std::thread::sleep;
use std::time::Duration;

// TODO : lots of error handling / chaos testing

fn get_interface() -> String {
    let args: Vec<String> = std::env::args().skip(1).collect();
    if args.len() == 0 {
        println!("\x1b[31merr:\x1b[0m -> need an interface argument");
        std::process::exit(1);
    }
    return args[0].to_owned();
}

fn parse_net_line(ln: &str) -> (u64, u64) {
    let mut rx: u64 = 0;
    let mut tx: u64 = 0;
    for (idx, val) in ln.split(" ").filter(|i| i.len() > 0).enumerate() {
        if idx == 1 {
            rx = val.parse().unwrap();
        } else if idx == 9 {
            tx = val.parse().unwrap();
        }
    }
    return (rx, tx);
}

// const RX: &'static str = "%{{F#696969}}﮴%{{F-}}";
// const TX: &'static str = "%{{F#696969}}﮵%{{F-}}";
// const NEUTRAL: &'static str = "%{{F#696969}}-%{{F-}}";
const RX: char = '﮴';
const TX: char = '﮵';
const NEUTRAL: char = '-';
// const DOWN: char = '睊';
const UP: char = '直';
const GREEN: &'static str = "#11F71D";
const MAGENTA: &'static str = "#AD003D";

fn print_status(more_rx: bool, more_tx: bool) {
    let mut r = NEUTRAL;
    let mut t = NEUTRAL;
    if more_rx {
        r = RX;
    }
    if more_tx {
        t = TX;
    }

    println!(
        "%{{F{}}}NET%{{F-}} %{{F{}}}{}%{{F-}}  {}{}",
        MAGENTA, GREEN, UP, r, t
    );
}

fn main() {
    let net_dev = get_interface();

    let mut buf: [u8; 1000] = [0; 1000];
    let fresult = File::options().read(true).open("/proc/net/dev");
    let net_file: File;
    match fresult {
        Ok(f) => net_file = f,
        Err(e) => {
            poly::polyerr(e);
            std::process::exit(1);
        }
    }

    let mut found = false;
    let mut n: usize;
    let mut prev_rx: u64 = 0;
    let mut prev_tx: u64 = 0;
    let mut rx: u64 = 0;
    let mut tx: u64 = 0;
    let mut more_rx = false;
    let mut more_tx = false;

    // seems to be where /proc/net/dev's header stops
    const OFFSET: u64 = 200;

    loop {
        let rdres = net_file.read_at(&mut buf, OFFSET);
        match rdres {
            Ok(_n) => n = _n,
            Err(e) => {
                poly::polyerr(e);
                std::process::exit(1);
            }
        }

        let traf = String::from_utf8_lossy(&buf[..n]);
        for line in traf.lines() {
            if line[..10].contains(&net_dev) {
                found = true;
                (rx, tx) = parse_net_line(line);
            }
        }

        if !found {
            poly::polyerr(format!("ERR: {} not found", net_dev));
            std::process::exit(1);
        }

        if rx > prev_rx {
            more_rx = true;
        } else {
            more_rx = false;
        }
        prev_rx = rx;

        if tx > prev_tx {
            more_tx = true;
        } else {
            more_tx = false;
        }
        prev_tx = tx;

        print_status(more_rx, more_tx);
        sleep(Duration::from_millis(500));
    }
}
