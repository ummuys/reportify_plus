// ---- старт ----
import { API_BASE } from './config/index.js';
import { loadSchemas } from './api/index.js';
import { syncRoleFromToken } from './core/index.js';
import { setupEventListeners, updateButtons, updateSchemaSelect, refreshHistory, showAlert, initChartModal, applyRoleRestrictions } from './ui/index.js';
import { loader } from './ui/loader.js';

function initTypedHeading() {
  const target = document.querySelector('.typed-text');
  if (!target) return;

  const TypedConstructor = window.Typed;
  if (typeof TypedConstructor !== 'function') {
    console.warn('Typed.js не найден на window');
    return;
  }

  if (target.dataset.typedReady === 'true') {
    return;
  }

  target.dataset.typedReady = 'true';

  new TypedConstructor(target, {
    strings: [
      'Конструктор отчёта',
      'Создавайте отчёты быстрее'
    ],
    typeSpeed: 45,
    backSpeed: 25,
    backDelay: 2200,
    smartBackspace: true,
    loop: true
  });
}

async function init() {
	const token = localStorage.getItem('access_token_v1') || ''
	console.log(
		'Используем токен для API:',
		token ? 'присутствует' : 'отсутствует',
	)
	const currentRole = syncRoleFromToken(token)
	applyRoleRestrictions(currentRole)

	initTypedHeading()

	if (!token) {
		loader.hide()
		await showAlert('Требуется авторизация. Войдите в систему.', 'Авторизация')
		window.location.assign('/')
		return
	}

	try {
		loader.show()

		console.log('Начинаем загрузку схем...')
		const schemasData = await loadSchemas()
		console.log('Получены данные схем:', schemasData)

		if (!schemasData || typeof schemasData !== 'object') {
			throw new Error('Некорректный формат данных схем')
		}

		const normalizedSchemas = {
			user: Array.isArray(schemasData.user) ? schemasData.user : [],
			sys: Array.isArray(schemasData.sys) ? schemasData.sys : [],
		}

		console.log('Нормализованные схемы:', normalizedSchemas)
		await updateSchemaSelect(normalizedSchemas)
		console.log('Схемы успешно загружены в select')

		setupEventListeners()

		// ✅ ИЗМЕНЕНО: История загружается без API кэша (только localStorage)
		// API кэша не существует, поэтому используем только локальное хранилище
		refreshHistory({ silent: true }).catch(err => {
			console.warn(
				'История загружается только из localStorage (API кэша недоступен)',
			)
		})

		updateButtons()
		initChartModal()
	} catch (e) {
		console.error('Ошибка при инициализации:', e)
		const schemaError = document.getElementById('schemaError')
		if (schemaError) {
			schemaError.textContent =
				'Не удалось загрузить схемы БД: ' + (e.message || e)
			schemaError.style.display = 'block'
		}
		await showAlert(
			'Не удалось загрузить схемы базы данных.\n' + (e.message || e),
			'Ошибка инициализации',
		)
	} finally {
		setTimeout(() => {
			loader.hide()
		}, 500)
	}
}


// Запускаем когда DOM полностью загружен
if (document.readyState === 'loading') {
  document.addEventListener('DOMContentLoaded', init);
} else {
  // DOM уже готов
  init();
}
