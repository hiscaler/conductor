import importlib.util, os
p=r'D:\wwwroot\hiscaler\conductor\output\Temu-US\批202610130955-月夜吸血鬼骷髅咖啡\图片\render_fixed_print.py'
spec=importlib.util.spec_from_file_location('r',p); m=importlib.util.module_from_spec(spec); spec.loader.exec_module(m)
for sku in ['MUG0110GN','MUG0110BL']:
 im=m.cutout_sku(sku); im.save(os.path.join(os.environ['TEMP'],sku+'-cut.png'))
 print(sku, im.size)
