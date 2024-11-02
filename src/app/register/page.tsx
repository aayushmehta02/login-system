"use client";

import { useState } from 'react';

export default function Register(): JSX.Element {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [empName, setEmpName] = useState('');
  const [empAge, setEmpAge] = useState('');
  const [department, setDepartment] = useState('');

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const response = await fetch('http://localhost:8000/api/register', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ username, password, empName, empAge: parseInt(empAge), department }),
    });
    const data = await response.json();
    console.log(data.message); // For feedback on registration
  };

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-100">
      <div className="w-full max-w-sm p-8 bg-white rounded-lg shadow-md">
        <h1 className="mb-6 text-2xl font-bold text-center">Register</h1>
        <form onSubmit={handleSubmit}>
          <div className="mb-4">
            <label htmlFor="username" className="block mb-2 text-sm font-medium text-gray-700">Username:</label>
            <input
              type="text"
              id="username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              required
              className="w-full p-2 border border-gray-300 rounded-md focus:outline-none focus:ring focus:ring-blue-300 text-black"
            />
          </div>
          <div className="mb-4">
            <label htmlFor="password" className="block mb-2 text-sm font-medium text-gray-700">Password:</label>
            <input
              type="password"
              id="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              className="w-full p-2 border border-gray-300 rounded-md focus:outline-none focus:ring focus:ring-blue-300 text-black"
            />
          </div>
          <div className="mb-4">
            <label htmlFor="empName" className="block mb-2 text-sm font-medium text-gray-700">Employee Name:</label>
            <input
              type="text"
              id="empName"
              value={empName}
              onChange={(e) => setEmpName(e.target.value)}
              required
              className="w-full p-2 border border-gray-300 rounded-md focus:outline-none focus:ring focus:ring-blue-300 text-black"
            />
          </div>
          <div className="mb-4">
            <label htmlFor="empAge" className="block mb-2 text-sm font-medium text-gray-700">Age:</label>
            <input
              type="number"
              id="empAge"
              value={empAge}
              onChange={(e) => setEmpAge(e.target.value)}
              required
              className="w-full p-2 border border-gray-300 rounded-md focus:outline-none focus:ring focus:ring-blue-300 text-black"
            />
          </div>
          <div className="mb-6">
            <label htmlFor="department" className="block mb-2 text-sm font-medium text-gray-700">Department:</label>
            <input
              type="text"
              id="department"
              value={department}
              onChange={(e) => setDepartment(e.target.value)}
              required
              className="w-full p-2 border border-gray-300 rounded-md focus:outline-none focus:ring focus:ring-blue-300 text-black"
            />
          </div>
          <button
            type="submit"
            className="w-full py-2 font-semibold text-white bg-blue-600 rounded-md hover:bg-blue-700 focus:outline-none focus:ring focus:ring-blue-300"
          >
            Register
          </button>
        </form>
      </div>
    </div>
  );
}
