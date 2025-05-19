// src/routes/+layout.server.ts
import { redirect } from '@sveltejs/kit';
import jwt from 'jsonwebtoken';
import type { LayoutServerLoad } from './$types';

type UserInfo = {
	email: string;
	employee_id: string;
	employee_name: string;
	primary_role: number;
	role: string;
};

export const load: LayoutServerLoad = async ({ cookies }) => {
	const token = cookies.get('token');

	if (!token) {
		throw redirect(302, '/login');
	}

	let user: UserInfo | null = null;

	try {
		const decoded = jwt.decode(token) as UserInfo | null;

		if (!decoded) {
			throw redirect(302, '/login');
		}

		user = {
			email: decoded.email,
			employee_id: decoded.employee_id,
			employee_name: decoded.employee_name,
			primary_role: decoded.primary_role,
			role: decoded.role
		};
	} catch (err) {
		console.error('JWT decode error:', err);
		throw redirect(302, '/login');
	}

	return { user };
};
