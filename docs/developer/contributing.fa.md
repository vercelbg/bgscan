<div align="left">

[**English**](./contributing.md)  |  [**فارسی**](./contributing.fa.md)

</div>


# راهنمای مشارکت (Contributing Guide)

از اینکه برای مشارکت در **bgscan** وقت می‌گذارید، متشکریم! هر نوع مشارکتی ارزشمند است؛ چه رفع Bug باشد، چه بهبود Documentation یا اضافه کردن Featureهای جدید.

## قبل از شروع

* قبل از ایجاد یک **Issue** جدید، Issueهای موجود را بررسی کنید.
* اگر قصد پیاده‌سازی یک Feature بزرگ یا ایجاد تغییرات اساسی را دارید، ابتدا یک **Issue** باز کنید تا درباره آن گفتگو شود.
* هر **Pull Request** فقط روی یک تغییر مشخص متمرکز باشد.

---

## Branch Strategy

لطفاً **مستقیماً روی Branch `main` Commit یا Push نکنید.**

ابتدا از `main` یک Branch جدید ایجاد کنید:

```bash
git checkout main
git pull origin main
git checkout -b feature/my-feature
```

نام‌های پیشنهادی برای Branch:

```text
feature/add-xray-parser
feature/ui-improvements
fix/memory-leak
fix/windows-build
docs/update-readme
refactor/scanner
test/unit-tests
```

پس از اتمام تغییرات:

1. تغییرات خود را **Commit** کنید.
2. Branch را **Push** کنید.
3. یک **Pull Request** به Branch `main` این پروژه ارسال کنید.
4. پس از **Merge** شدن Pull Request، Branch خود را حذف کنید.

> **نکته برای Maintainer:** Maintainerهای پروژه می‌توانند برای تغییرات کوچک مستقیماً روی Branch `main` Commit و Push کنند. برای تغییرات بزرگ، بهتر است ابتدا یک Branch جداگانه ایجاد کرده و پس از بررسی، آن را در `main` Merge کنند.

---

## Coding Style

* از Style و Conventionهای فعلی پروژه پیروی کنید.
* Code را تا حد امکان ساده و خوانا نگه دارید.
* Commitهای کوچک و متمرکز ایجاد کنید.
* از اعمال تغییرات Formatting غیرمرتبط خودداری کنید.
* فقط در صورت نیاز Comment اضافه کنید.

---

## Commit Messages

از Commit Messageهای واضح و توصیفی استفاده کنید.

نمونه‌ها:

```text
feat: add HTTP/3 support
fix: resolve race condition in scanner
docs: improve installation guide
refactor: simplify writer pipeline
test: add unit tests for parser
```

---

## Pull Requests

قبل از ارسال **Pull Request** مطمئن شوید که:

* پروژه بدون خطا **Build** می‌شود.
* تمام Testهای موجود با موفقیت اجرا می‌شوند.
* قابلیت جدید به‌درستی Test شده است.
* در صورت نیاز، Documentation نیز به‌روزرسانی شده است.

لطفاً در Pull Request موارد زیر را توضیح دهید:

* چه چیزی تغییر کرده است.
* دلیل این تغییر چیست.
* Screenshot (در صورت وجود تغییرات UI)
* Issue مرتبط (در صورت وجود)

---

## Reporting Bugs

هنگام گزارش Bug، تا حد امکان اطلاعات زیر را ارائه دهید:

* Operating System
* Architecture
* نسخه bgscan
* Configuration
* مراحل بازتولید مشکل
* رفتار مورد انتظار
* رفتار فعلی
* Logs (در صورت امکان)

---

## Feature Requests

در صورت پیشنهاد یک Feature جدید، لطفاً موارد زیر را توضیح دهید:

* چه مشکلی را می‌خواهید حل کنید.
* راهکار پیشنهادی شما چیست.
* آیا راهکار جایگزینی وجود دارد؟
* هرگونه توضیح یا Context اضافی.

---

## Code of Conduct

با سایر مشارکت‌کنندگان با احترام و روحیه همکاری برخورد کنید.

بحث‌های فنی و سازنده کاملاً مورد استقبال هستند، اما توهین، آزار، رفتار غیرمحترمانه یا حملات شخصی در این پروژه پذیرفته نخواهد شد.

از اینکه به بهتر شدن **bgscan** کمک می‌کنید، سپاسگزاریم.
