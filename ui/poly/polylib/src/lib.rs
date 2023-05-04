use std::fmt::Display;

pub fn polyfmt<T>(s: T, colorhex: T)
where
    T: Display,
{
    println!("%{{F{}}}{}%{{F-}}", colorhex, s);
}

// Red hex. Not 100% sure when it needs newline terminator
// if something crashes
pub fn polyerr<T>(err: T)
where
    T: Display,
{
    println!("%{{F#F54242}}{}%{{F-}}", err);
}

/// To test hex colors in equivalent rgb. Prints to stdout.
pub fn debug_print<T>(s: T, r: i32, g: i32, b: i32)
where
    T: Display,
{
    println!("\x1b[38;2;{};{};{}m{}\x1b[0m", r, g, b, s);
}
