INSERT INTO services
(category, code, name, price, program_id, sid, renewal_day, trial_day, url_telco, url_portal, url_callback, url_notif_sub, url_notif_unsub, url_notif_renewal, url_postback)
VALUES
('CLOUDPLAY', 'CLOUDPLAY', 'CLOUDPLAY', 2220, 'CLOUDPLAY', 'GAMGENRLITCLOUDPLAY_Subs', 2, 0, 'https://api.digitalcore.telkomsel.com', 'https://id.cloudplay.mobi', 'https://id.cloudplay.mobi/subscription/login', 'https://id.cloudplay.mobi/api/subscription/subscribe', 'https://id.cloudplay.mobi/api/subscription/unsubscribe', 'https://id.cloudplay.mobi/api/subscription/renewal', 'http://kbtools.net/id-linkittisel.php'),
('CLOUDPLAY', 'CLOUDPLAY1', 'CLOUDPLAY 1', 3330, 'CLOUDPLAY1', 'GAMGENRLITCLOUDPLAY1_Subs', 3, 0, 'https://api.digitalcore.telkomsel.com', 'https://id.cloudplay.mobi', 'https://id.cloudplay.mobi/subscription/login', 'https://id.cloudplay.mobi/api/subscription/subscribe', 'https://id.cloudplay.mobi/api/subscription/unsubscribe', 'https://id.cloudplay.mobi/api/subscription/renewal', 'http://kbtools.net/id-linkittisel.php'),
('CLOUDPLAY', 'CLOUDPLAY2', 'CLOUDPLAY 2', 5550, 'CLOUDPLAY2', 'GAMGENRLITCLOUDPLAY2_Subs', 7, 0, 'https://api.digitalcore.telkomsel.com', 'https://id.cloudplay.mobi', 'https://id.cloudplay.mobi/subscription/login', 'https://id.cloudplay.mobi/api/subscription/subscribe', 'https://id.cloudplay.mobi/api/subscription/unsubscribe', 'https://id.cloudplay.mobi/api/subscription/renewal', 'http://kbtools.net/id-linkittisel.php'),
('CLOUDPLAY', 'CLOUDPLAY3', 'CLOUDPLAY 3', 11100, 'CLOUDPLAY3', 'GAMGENRLITCLOUDPLAY3_Subs', 14, 0, 'https://api.digitalcore.telkomsel.com', 'https://id.cloudplay.mobi', 'https://id.cloudplay.mobi/subscription/login', 'https://id.cloudplay.mobi/api/subscription/subscribe', 'https://id.cloudplay.mobi/api/subscription/unsubscribe', 'https://id.cloudplay.mobi/api/subscription/renewal', 'http://kbtools.net/id-linkittisel.php'),
('CLOUDPLAY', 'CLOUDPLAY4', 'CLOUDPLAY 4', 16650, 'CLOUDPLAY4', 'GAMGENRLITCLOUDPLAY4_Subs', 30, 0, 'https://api.digitalcore.telkomsel.com', 'https://id.cloudplay.mobi', 'https://id.cloudplay.mobi/subscription/login', 'https://id.cloudplay.mobi/api/subscription/subscribe', 'https://id.cloudplay.mobi/api/subscription/unsubscribe', 'https://id.cloudplay.mobi/api/subscription/renewal', 'http://kbtools.net/id-linkittisel.php');


INSERT INTO contents
(service_id, name, value, tid)
VALUES
(1, 'FIRSTPUSH', '[2220] Cloudplay Games bisa kamu mainkan sekarang klik id.cloudplay.mobi (berlaku tarif internet) PIN: @pin Stop: UNREG CLOUDPLAY ke 97770 CS:02152902182', '69'),
(1, 'RENEWAL', '[2220] Cloudplay Games bisa kamu mainkan sekarang klik id.cloudplay.mobi (berlaku tarif internet) PIN: @pin Stop: UNREG CLOUDPLAY ke 97770 CS:02152902182', '69'),
(2, 'FIRSTPUSH', '[3330] Cloudplay Games bisa kamu mainkan sekarang klik id.cloudplay.mobi (berlaku tarif internet) PIN: @pin Stop: UNREG CLOUDPLAY ke 97770 CS:02152902182', '96'),
(2, 'RENEWAL', '[3330] Cloudplay Games bisa kamu mainkan sekarang klik id.cloudplay.mobi (berlaku tarif internet) PIN: @pin Stop: UNREG CLOUDPLAY ke 97770 CS:02152902182', '96'),
(3, 'FIRSTPUSH', '[5550] Cloudplay Games bisa kamu mainkan sekarang klik id.cloudplay.mobi (berlaku tarif internet) PIN: @pin Stop: UNREG CLOUDPLAY ke 97770 CS:02152902182', '142'),
(3, 'RENEWAL', '[5550] Cloudplay Games bisa kamu mainkan sekarang klik id.cloudplay.mobi (berlaku tarif internet) PIN: @pin Stop: UNREG CLOUDPLAY ke 97770 CS:02152902182', '142'),
(4, 'FIRSTPUSH', '[11100] Cloudplay Games bisa kamu mainkan sekarang klik id.cloudplay.mobi (berlaku tarif internet) PIN: @pin Stop: UNREG CLOUDPLAY ke 97770 CS:02152902182', '43'),
(4, 'RENEWAL', '[11100] Cloudplay Games bisa kamu mainkan sekarang klik id.cloudplay.mobi (berlaku tarif internet) PIN: @pin Stop: UNREG CLOUDPLAY ke 97770 CS:02152902182', '43'),
(5, 'FIRSTPUSH', '[16650] Cloudplay Games bisa kamu mainkan sekarang klik id.cloudplay.mobi (berlaku tarif internet) PIN: @pin Stop: UNREG CLOUDPLAY ke 97770 CS:02152902182', '58'),
(5, 'RENEWAL', '[16650] Cloudplay Games bisa kamu mainkan sekarang klik id.cloudplay.mobi (berlaku tarif internet) PIN: @pin Stop: UNREG CLOUDPLAY ke 97770 CS:02152902182', '58');

INSERT INTO schedules
(id, name, publish_at, unlocked_at, is_unlocked)
VALUES
(1, 'REMINDER', NOW(), NOW(), false),
(2, 'RENEWAL', NOW(), NOW(), false),
(3, 'RETRY', NOW(), NOW(), false);

INSERT INTO adnets
(name, value)
VALUES
('adn', 'adn1');
