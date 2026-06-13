api: sh -c "cleanup() { just _compose down $SERVICES; }; trap cleanup EXIT; just _compose up --build --detach $SERVICES && just _compose logs --follow $SERVICES"
front-identityx: cd front/identityx && pnpm dev
front-informd: cd front/informd && pnpm dev
front-payssage: cd front/payssage && pnpm dev
front-univents: cd front/univents && pnpm dev
