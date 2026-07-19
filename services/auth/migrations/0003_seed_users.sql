-- Demo host accounts with fixed UUIDs so the main API seed can reference them.
INSERT INTO users (id, phone, full_name, avatar_url, is_host, status) VALUES
    ('11111111-1111-1111-1111-111111111101', '+989121110001', 'آرش کریمی',    'https://i.pravatar.cc/200?u=arash-karimi',    true, 'active'),
    ('11111111-1111-1111-1111-111111111102', '+989121110002', 'مریم حسینی',   'https://i.pravatar.cc/200?u=maryam-hosseini', true, 'active'),
    ('11111111-1111-1111-1111-111111111103', '+989121110003', 'رضا صادقی',    'https://i.pravatar.cc/200?u=reza-sadeghi',    true, 'active'),
    ('11111111-1111-1111-1111-111111111104', '+989121110004', 'نگار احمدی',   'https://i.pravatar.cc/200?u=negar-ahmadi',    true, 'active'),
    ('11111111-1111-1111-1111-111111111105', '+989121110005', 'سینا مرادی',   'https://i.pravatar.cc/200?u=sina-moradi',     true, 'active'),
    ('11111111-1111-1111-1111-111111111106', '+989121110006', 'فاطمه رستمی',  'https://i.pravatar.cc/200?u=fatemeh-rostami', true, 'active');
