<div align="left">

[**English**](./README.md)  |  [**فارسی**](./README.fa.md)

</div>

# راهنمای توسعه‌دهنده

خوش آمدید! این راهنما نحوه راه‌اندازی محیط، بیلد و اجرای پروژه **bgscan** را به‌صورت محلی توضیح می‌دهد.

> **قبل از شروع:** لطفاً ابتدا [راهنمای مشارکت (Contributing Guide)](./contributing.fa.md) را مطالعه کنید.

## فهرست مطالب

- [ساختار پروژه](#ساختار-پروژه)
- [پیش‌نیازها](#پیش‌نیازها)
- [شروع کار](#شروع-کار)
  - [۱. کلون کردن مخزن](#۱-کلون-کردن-مخزن)
  - [۲. ساخت یک برنچ](#۲-ساخت-یک-برنچ)
  - [۳. نصب وابستگی‌ها](#۳-نصب-وابستگی‌ها)
  - [۴. اجرای پروژه](#۴-اجرای-پروژه)
- [بیلد گرفتن نسخه‌های انتشار (Release)](#بیلد-گرفتن-نسخه‌های-انتشار-release)
  - [نصب builder](#نصب-builder)
  - [دستورات بیلد](#دستورات-بیلد)
  - [بیلد گرفتن برای اندروید](#بیلد-گرفتن-برای-اندروید)
- [مرجع bgscan-builder](#مرجع-bgscan-builder)

---

## ساختار پروژه

```
.
├── assets
│   └── xray
│       └── outbounds
├── bgscan-builder
├── build.sh
├── cmd
│   └── bgscan
│       └── main.go
├── docs
│   ├── developer
│   ├── images
│   └── user
├── go.mod
├── images
├── internal
│   ├── core
│   │   ├── config
│   │   ├── dns
│   │   ├── fileutil
│   │   ├── ip
│   │   ├── iplist
│   │   ├── process
│   │   ├── result
│   │   ├── scanner
│   │   └── xray
│   ├── logger
│   ├── startup
│   └── ui
│       ├── components
│       ├── main
│       ├── shared
│       └── theme
├── ips
├── LICENSE
├── README.md
├── scripts
│   ├── install-builder.ps1
│   ├── install-builder.sh
│   ├── install-deps.ps1
│   ├── install-deps.sh
│   ├── install.ps1
│   ├── install.sh
│   ├── publish.sh
│   └── release.sh
└── settings
```

نگاهی سریع به پوشه‌های اصلی:

| مسیر | توضیح |
|---|---|
| `cmd/bgscan` | نقطه ورود برنامه (`main.go`). |
| `internal/core` | منطق اصلی: تنظیمات (config)، DNS، ابزارهای فایل، مدیریت IP، مدیریت پروسه، نتایج، اسکنر و یکپارچه‌سازی با Xray. |
| `internal/logger` | ابزارهای لاگ‌گیری. |
| `internal/startup` | بررسی‌های سلامت در زمان راه‌اندازی (config، DNS، logger، Xray). |
| `internal/ui` | کامپوننت‌های UI، نمای اصلی، عناصر مشترک و تم. |
| `assets/xray/outbounds` | فایل‌های outbound تعبیه‌شده Xray. |
| `ips` | لیست‌های پیش‌فرض رنج IP به ازای هر ارائه‌دهنده (Cloudflare، AWS، Azure و غیره). |
| `settings` | فایل‌های تنظیمات پیش‌فرض با فرمت `.toml`. |
| `scripts` | اسکریپت‌های کمکی برای نصب، بیلد و انتشار. |
| `docs` | مستندات کاربر و توسعه‌دهنده. |

---

## پیش‌نیازها

- [Go](https://go.dev/) (نسخه مورد نیاز را در `go.mod` ببینید)
- Git
- برای بیلد اندروید: [Android NDK](https://developer.android.com/ndk)

---

## شروع کار

### ۱. کلون کردن مخزن

```bash
git clone <repo-url>
cd bgscan
```

### ۲. ساخت یک برنچ

```bash
git checkout -b feature/my-change
```

### ۳. نصب وابستگی‌ها

bgscan از ابزاری کمکی به نام **`bgscan-builder`** برای دریافت و ساخت وابستگی‌های پروژه استفاده می‌کند. نیازی نیست آن را به‌صورت دستی نصب کنید — اسکریپت‌های نصب زیر آن را برای شما دانلود کرده، در مسیر اصلی پروژه قرار می‌دهند و از آن برای دریافت نسخه صحیح وابستگی متناسب با سیستم‌عامل و معماری شما استفاده می‌کنند.

از مسیر اصلی پروژه، دستور زیر را اجرا کنید:

**لینوکس / مک‌اواس**
```bash
./scripts/install-deps.sh
```

**ویندوز**
```powershell
./scripts/install-deps.ps1
```

این اسکریپت موارد زیر را انجام می‌دهد:
۱. دانلود `bgscan-builder` و قرار دادن آن در مسیر اصلی پروژه.
۲. اجرای `bgscan-builder setup-dev --project-dir <project-root>` برای دانلود وابستگی‌های صحیح متناسب با سیستم‌عامل/معماری شما و قرار دادن آن‌ها در مسیر درست.

### ۴. اجرای پروژه

پس از نصب وابستگی‌ها:

```bash
go mod tidy
go run ./cmd/bgscan/
```

---

## بیلد گرفتن نسخه‌های انتشار (Release)

برای ساخت فایل‌های انتشار (release)، به `bgscan-builder` نیز نیاز دارید. اگر آن را در مرحله نصب وابستگی‌ها دریافت نکرده‌اید، با دستور زیر نصبش کنید:

**لینوکس / مک‌اواس**
```bash
./scripts/install-builder.sh
```

**ویندوز**
```powershell
./scripts/install-builder.ps1
```

### دستورات بیلد

```bash
bgscan-builder release -os linux -arch amd64
bgscan-builder release -os android -arch arm64 -ndk-dir /opt/android-ndk
bgscan-builder release -os all -arch all -dest ./dist
```

### بیلد گرفتن برای اندروید

بیلد اندروید نیاز به Android NDK دارد. مسیر آن را با `-ndk-dir` مشخص کنید:

```bash
bgscan-builder release -os android -arch arm64 -ndk-dir /opt/android-ndk
```

---

## مرجع bgscan-builder

فایل‌های انتشار را برای یک یا چند ترکیب سیستم‌عامل/معماری بیلد می‌کند.

**مثال‌ها:**

```bash
bgscan-builder release -os linux -arch amd64
bgscan-builder release -os android -arch arm64 -ndk-dir /opt/android-ndk
bgscan-builder release -os all -arch all -dest ./dist
```

**فلگ‌ها (Flags):**

| فلگ | توضیح |
|---|---|
| `-arch string` | معماری هدف (`amd64`, `arm64`, `arm32`, `amd32`, `all`) |
| `-dep-version string` | تگ نسخه وابستگی‌ها (پیش‌فرض `"v1.0"`) |
| `-dest string` | مسیر خروجی انتشار (پیش‌فرض `"./dist"`) |
| `-ndk-dir string` | مسیر ریشه Android NDK |
| `-os string` | سیستم‌عامل هدف (`linux`, `windows`, `macos`, `android`, `all`) |
| `-project-dir string` | مسیر پروژه bgscan |
| `-xray-version string` | تگ نسخه Xray (پیش‌فرض `"v26.3.27"`) |
