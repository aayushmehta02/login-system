"use client"

import { useEffect, useState } from 'react';

export default function Register(): JSX.Element {
  const [formData, setFormData] = useState({
    username: '',
    password: '',
    firstName: '',
    middleName: '',
    lastName: '',
    contactNo: '',
    email: '',
    dob: '',
    address: '',
    state: '',
    city: '',
    pin: '',
    aadhar: '',
    pan: '',
    active: true,
  });
  const [serverMessage, setServerMessage] = useState('');
  
  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value, type } = e.target;
  
    // Check for the 'checked' property if the input type is 'checkbox'
    const isCheckbox = type === 'checkbox';
  
    setFormData({
      ...formData,
      [name]: isCheckbox ? (e.target as HTMLInputElement).checked : value,
    });
  };

  useEffect(() => {
    // Example: Initialize form data or perform some side effect
    console.log('Component mounted or formData updated');
  }, [formData]); // Runs when formData changes

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const formattedDob = formData.dob ? new Date(formData.dob).toISOString() : '';
    try {
      const response = await fetch('http://localhost:8000/api/register', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ ...formData, dob: formattedDob, }),
      });
      console.log(JSON.stringify({ ...formData, dob: formattedDob, }, null, 2));

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || 'Registration failed. Please try again.');
      }
  
      const data = await response.json();
      setServerMessage(data.message); // Save the server response to state
      alert(data.message); // Feedback on successful registration
    } catch (error: any) {
      setServerMessage(error.message || 'An unexpected error occurred. Please try again.');
      alert(`Error: ${error.message || 'An unexpected error occurred. Please try again.'}`);
    }
  };
  

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-100">
      <div className="w-full max-w-md p-8 bg-white rounded-lg shadow-md">
        <h1 className="mb-6 text-2xl font-bold text-center text-gray-700">Register</h1>
        <form onSubmit={handleSubmit}>
          {/** Username and Password */}
          <div className="mb-4">
            <label htmlFor="username" className="block mb-2 text-sm font-medium text-gray-700">Username:</label>
            <input
              type="text"
              name="username"
              id="username"
              value={formData.username}
              onChange={handleChange}
              required
              className="w-full p-2 border border-gray-300 rounded-md focus:outline-none focus:ring focus:ring-blue-300"
            />
          </div>
          <div className="mb-4">
            <label htmlFor="password" className="block mb-2 text-sm font-medium text-gray-700">Password:</label>
            <input
              type="password"
              name="password"
              id="password"
              value={formData.password}
              onChange={handleChange}
              required
              className="w-full p-2 border border-gray-300 rounded-md focus:outline-none focus:ring focus:ring-blue-300"
            />
          </div>

          <div className="mb-4">
  <label htmlFor="firstName" className="block mb-2 text-sm font-medium text-gray-700">
    First Name:
  </label>
  <input
    type="text"
    name="firstName"
    id="firstName"
    value={formData.firstName}
    onChange={handleChange}
    className="w-full p-2 border border-gray-300 rounded-md focus:outline-none focus:ring focus:ring-blue-300"
  />
</div>

<div className="mb-4">
  <label htmlFor="middleName" className="block mb-2 text-sm font-medium text-gray-700">
    Middle Name:
  </label>
  <input
    type="text"
    name="middleName"
    id="middleName"
    value={formData.middleName}
    onChange={handleChange}
    className="w-full p-2 border border-gray-300 rounded-md focus:outline-none focus:ring focus:ring-blue-300"
  />
</div>

<div className="mb-4">
  <label htmlFor="lastName" className="block mb-2 text-sm font-medium text-gray-700">
    Last Name:
  </label>
  <input
    type="text"
    name="lastName"
    id="lastName"
    value={formData.lastName}
    onChange={handleChange}
    className="w-full p-2 border border-gray-300 rounded-md focus:outline-none focus:ring focus:ring-blue-300"
  />
</div>


          {/** Contact Information */}
          <div className="mb-4">
            <label htmlFor="contactNo" className="block mb-2 text-sm font-medium text-gray-700">Contact No:</label>
            <input
              type="text"
              name="contactNo"
              id="contactNo"
              value={formData.contactNo}
              onChange={handleChange}
              className="w-full p-2 border border-gray-300 rounded-md focus:outline-none focus:ring focus:ring-blue-300"
            />
          </div>
          <div className="mb-4">
            <label htmlFor="email" className="block mb-2 text-sm font-medium text-gray-700">Email:</label>
            <input
              type="email"
              name="email"
              id="email"
              value={formData.email}
              onChange={handleChange}
              className="w-full p-2 border border-gray-300 rounded-md focus:outline-none focus:ring focus:ring-blue-300"
            />
          </div>

          {/** Additional Fields */}
          <div className="mb-4">
            <label htmlFor="dob" className="block mb-2 text-sm font-medium text-gray-700">Date of Birth:</label>
            <input
              type="date"
              name="dob"
              id="dob"
              value={formData.dob}
              onChange={handleChange}
              className="w-full p-2 border border-gray-300 rounded-md focus:outline-none focus:ring focus:ring-blue-300"
            />
          </div>
          <div className="mb-4">
            <label htmlFor="address" className="block mb-2 text-sm font-medium text-gray-700">Address:</label>
            <input
              type="text"
              name="address"
              id="address"
              value={formData.address}
              onChange={handleChange}
              className="w-full p-2 border border-gray-300 rounded-md focus:outline-none focus:ring focus:ring-blue-300"
            />
          </div>
          <div className="grid grid-cols-2 gap-4 mb-4">
            <div>
              <label htmlFor="state" className="block mb-2 text-sm font-medium text-gray-700">State:</label>
              <input
                type="text"
                name="state"
                id="state"
                value={formData.state}
                onChange={handleChange}
                className="w-full p-2 border border-gray-300 rounded-md focus:outline-none focus:ring focus:ring-blue-300"
              />
            </div>
            <div>
              <label htmlFor="city" className="block mb-2 text-sm font-medium text-gray-700">City:</label>
              <input
                type="text"
                name="city"
                id="city"
                value={formData.city}
                onChange={handleChange}
                className="w-full p-2 border border-gray-300 rounded-md focus:outline-none focus:ring focus:ring-blue-300"
              />
            </div>
          </div>

          {/** Aadhar, PAN, and Pin */}
          <div className="mb-4">
            <label htmlFor="pin" className="block mb-2 text-sm font-medium text-gray-700">Pin:</label>
            <input
              type="text"
              name="pin"
              id="pin"
              value={formData.pin}
              onChange={handleChange}
              className="w-full p-2 border border-gray-300 rounded-md focus:outline-none focus:ring focus:ring-blue-300"
            />
          </div>
          <div className="mb-4">
            <label htmlFor="aadhar" className="block mb-2 text-sm font-medium text-gray-700">Aadhar:</label>
            <input
              type="text"
              name="aadhar"
              id="aadhar"
              value={formData.aadhar}
              onChange={handleChange}
              className="w-full p-2 border border-gray-300 rounded-md focus:outline-none focus:ring focus:ring-blue-300"
            />
          </div>
          <div className="mb-4">
            <label htmlFor="pan" className="block mb-2 text-sm font-medium text-gray-700">PAN:</label>
            <input
              type="text"
              name="pan"
              id="pan"
              value={formData.pan}
              onChange={handleChange}
              className="w-full p-2 border border-gray-300 rounded-md focus:outline-none focus:ring focus:ring-blue-300"
            />
          </div>

          {/** Active */}
          <div className="mb-6">
            <label className="flex items-center">
              <input
                type="checkbox"
                name="active"
                checked={formData.active}
                onChange={handleChange}
                className="mr-2 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
              />
              Active
            </label>
          </div>

          {/** Submit Button */}
          <button
            type="submit"
            className="w-full py-2 text-white bg-blue-600 rounded-md hover:bg-blue-700 focus:outline-none focus:ring focus:ring-blue-300"
          >
            Register
          </button>
        </form>
      </div>
    </div>
  );
}
